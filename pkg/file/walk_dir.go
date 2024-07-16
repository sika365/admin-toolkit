package file

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

// WalkDir recursively walks the directory up to the specified depth.
func WalkDir(root string, maxDepth int, reContentType *regexp.Regexp) (files MapFiles, err error) {
	files = NewMapFiles()

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the current depth
		currentDepth := len(strings.Split(filepath.ToSlash(path), "/")) - len(strings.Split(filepath.ToSlash(root), "/"))

		if currentDepth <= maxDepth && !d.IsDir() {
			if file := NewFile(reContentType, path); file.Error() == nil {
				files.Add(file)
			} else {
				// let continue
				return nil
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
