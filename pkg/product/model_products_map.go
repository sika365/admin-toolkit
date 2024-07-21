package product

type MapProducts map[string]*Product

func NewMapProducts(products ...*Product) MapProducts {
	ps := make(MapProducts)
	for _, p := range products {
		ps.Add(p)
	}
	return ps
}

func (mps MapProducts) Get(key string) *Product {
	return mps[key]
}

func (mps MapProducts) Add(prd *Product) MapProducts {
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

func (mps MapProducts) GetValues() Products {
	prds := make(Products, 0, len(mps))
	for _, p := range mps {
		prds = append(prds, p)
	}
	return prds
}
