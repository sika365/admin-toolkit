package image

type MapImages map[string]*Image

func NewMapImages(files ...*Image) MapImages {
	fs := make(MapImages)
	for _, f := range files {
		fs.Add(f)
	}
	return fs
}

func (fs MapImages) Get(key string) *Image {
	return fs[key]
}

func (fs MapImages) Add(file *Image) MapImages {
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

func (fs MapImages) GetValues() Images {
	files := make(Images, 0, len(fs))
	for _, f := range fs {
		files = append(files, f)
	}
	return files
}
