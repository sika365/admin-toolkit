package file

import (
	"regexp"
)

type MapFiles map[string]*File

func NewMapFiles(files ...*File) MapFiles {
	m := make(MapFiles)
	for _, f := range files {
		m[f.Hash] = f
	}
	return m
}

func (fs MapFiles) Get(key string) *File {
	return fs[key]
}

func (fs MapFiles) Add(file *File) MapFiles {
	if file.Hash == "" {
		return fs
	}

	fs[file.Hash] = file
	return fs
}

func (fs MapFiles) FromFiles(contentType *regexp.Regexp, paths ...string) MapFiles {
	for _, path := range paths {
		if f := NewFile(contentType, path); f.err != nil {
			return fs
		} else {
			fs[f.Hash] = f
			return fs
		}
	}
	return fs
}

func (fs MapFiles) GetKeys() []string {
	hashes := make([]string, 0, len(fs))
	for k := range fs {
		hashes = append(hashes, k)
	}
	return hashes
}

func (fs MapFiles) GetValues() Files {
	files := make(Files, 0, len(fs))
	for _, f := range fs {
		files = append(files, f)
	}
	return files
}
