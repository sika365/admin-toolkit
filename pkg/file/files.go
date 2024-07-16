package file

import (
	"regexp"
)

type Files []*File

func (fs Files) Add(file *File) Files {
	fs = append(fs, file)
	return fs
}

func (fs Files) AddFiles(contentType *regexp.Regexp, paths ...string) Files {
	for _, path := range paths {
		if f := NewFile(contentType, path); f.err != nil {
			return fs
		} else {
			fs = append(fs, f)
			return fs
		}
	}
	return fs
}

func (fs Files) GetHashes() []string {
	hashes := make([]string, len(fs))
	for i, f := range fs {
		hashes[i] = f.Hash
	}
	return hashes
}

func (fs Files) GetHashMap() map[string]*File {
	hashes := make(map[string]*File)
	for _, f := range fs {
		hashes[f.Hash] = f
	}
	return hashes
}

func (fs Files) GetHashMapKeys() ([]string, map[string]*File) {
	hashes := make([]string, len(fs))
	hashMap := make(map[string]*File)
	for i, f := range fs {
		hashes[i] = f.Hash
		hashMap[f.Hash] = f
	}
	return hashes, hashMap
}
