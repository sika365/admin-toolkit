package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// WalkDir recursively walks the directory up to the specified depth.
func WalkDir(root string, maxDepth int, reContentType *regexp.Regexp, reName *regexp.Regexp) (files MapFiles, err error) {
	files = NewMapFiles()

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the current depth
		currentDepth := len(strings.Split(filepath.ToSlash(path), "/")) - len(strings.Split(filepath.ToSlash(root), "/"))

		if currentDepth <= maxDepth {
			if !d.IsDir() {
				if file := NewFile(path, reContentType, reName); file.Error() == nil {
					files.Add(file)
				} else {
					// let continue
					return nil
				}
			} else if reName.MatchString(d.Name()) {
				// Read cover and gallery files
				// Read the directory
				imgFiles, err := os.ReadDir(path)
				if err != nil {
					logrus.Errorln("Error reading directory:", err)
					return err
				}

				// Loop through the directory entries
				for _, file := range imgFiles {
					if !file.IsDir() {
						f := &File{}

						if info, err := file.Info(); err != nil {
							f.err = err
							continue
						} else {
							f.Title = d.Name()
							f.Name = info.Name()
							f.Size = info.Size()
							f.ModTime = info.ModTime()
							f.Path = filepath.Join(path, f.Name)
						}

						if err := f.Open().
							MatchContentTypeFile(reContentType).
							GenerateHash().
							Close().err; err != nil {
							f.err = err
							continue
						} else {
							files.Add(f)
						}
					}
				}

				return filepath.SkipDir
			}
		}

		// Skip walking into subdirectories if maxDepth is reached
		if d.IsDir() && currentDepth >= maxDepth {
			return filepath.SkipDir
		}

		return nil
	})

	return files, err
}
