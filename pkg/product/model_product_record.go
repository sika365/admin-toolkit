package product

import (
	"github.com/sika365/admin-tools/pkg/category"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gorm.io/gorm"
)

type ProductRecords []*ProductRecord

type ProductRecord struct {
	models.CommonTableFields
	Barcode        string                  `json:"barcode,omitempty" sim:"primaryKey"`
	CategorySlug   string                  `json:"category,omitempty" sim:"primaryKey"`
	Title          string                  `json:"title,omitempty"`
	LocalProductID database.PID            `json:"local_product_id,omitempty" sim:"primaryKey"`
	LocalCategory  *category.LocalCategory `json:"local_category,omitempty" gorm:"foreignKey:CategorySlug;references:Slug"`
	LocalProduct   *LocalProduct           `json:"local_product,omitempty"`
}

func (pr *ProductRecord) BeforeCreate(tx *gorm.DB) error {
	if lp := pr.LocalProduct; lp == nil {
	} else if rlp := lp.Product; rlp == nil {
	} else if len(rlp.Tags) > 0 {
		rlp.Tags = nil
	}
	return nil
}

func (pr *ProductRecord) AddTopNodes(topNodes models.Nodes, replace bool) error {
	return pr.LocalProduct.AddTopNodes(topNodes, replace)
}

func (pr *ProductRecord) MatchBarcode(barcode string, prods models.Products) (product *models.Product, err error) {
	var products models.Products

	for _, p := range prods {
		if p.LocalProduct == nil {
			continue
		}
		for _, b := range p.LocalProduct.Barcodes {
			if b.Barcode == barcode {
				products = append(products, p)
				break
			}
		}
	}

	if len(products) == 0 {
		return nil, models.ErrNotFound
	} else {
		var prod = products[0]
		for _, p := range products {
			found := false
			for _, ps := range prod.LocalProduct.ProductStocks {
				if ps.StockID == p.ProductStock.StockID {
					found = true
					break
				}
			}
			if !found {
				prod.LocalProduct.ProductStocks = append(prod.LocalProduct.ProductStocks, &p.ProductStock)
			}
		}
		return prod, nil
	}
}
