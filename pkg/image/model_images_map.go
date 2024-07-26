package image

type MapImages map[string]*LocalImage

func NewMapImages(files ...*LocalImage) MapImages {
	fs := make(MapImages)
	for _, f := range files {
		fs.Add(f)
	}
	return fs
}

func (fs MapImages) Get(key string) *LocalImage {
	return fs[key]
}

func (fs MapImages) Add(file *LocalImage) MapImages {
	if h := file.Hash(); h == "" {
		return fs
	} else {
		fs[h] = file
		return fs
	}
}

func (fs MapImages) GetKeys() []string {
	hashes := make([]string, 0, len(fs))
	for k := range fs {
		hashes = append(hashes, k)
	}
	return hashes
}

func (fs MapImages) GetValues() LocalImages {
	files := make(LocalImages, 0, len(fs))
	for _, f := range fs {
		files = append(files, f)
	}
	return files
}
