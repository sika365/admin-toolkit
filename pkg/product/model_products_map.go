package product

type MapProducts map[string]*LocalProduct

func NewMapProducts(products ...*LocalProduct) MapProducts {
	ps := make(MapProducts)
	for _, p := range products {
		ps.Add(p)
	}
	return ps
}

func (mps MapProducts) Get(key string) *LocalProduct {
	return mps[key]
}

func (mps MapProducts) Add(prd *LocalProduct) MapProducts {
	if h := prd.Key(); h == "" {
		return mps
	} else {
		mps[h] = prd
		return mps
	}
}

func (mps MapProducts) GetKeys() []string {
	keys := make([]string, 0, len(mps))
	for k := range mps {
		keys = append(keys, k)
	}
	return keys
}

func (mps MapProducts) GetValues() LocalProducts {
	prds := make(LocalProducts, 0, len(mps))
	for _, p := range mps {
		prds = append(prds, p)
	}
	return prds
}
