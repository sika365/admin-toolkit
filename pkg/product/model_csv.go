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
	Barcode        string                  `json:"barcode,omitempty" gorm:"primaryKey"`
	Title          string                  `json:"title,omitempty"`
	CategoryAlias  string                  `json:"category,omitempty"`
	LocalProductID database.PID            `json:"local_product_id,omitempty"`
	LocalCategory  *category.LocalCategory `json:"local_category,omitempty" gorm:"foreignKey:CategoryAlias;references:Alias"`
	LocalProduct   *LocalProduct           `json:"local_product,omitempty"`
}

func (pr *ProductRecord) BeforeCreate(tx *gorm.DB) error {
	if lp := pr.LocalProduct; lp == nil {
	} else if rp := lp.Product; rp == nil {
	} else if rlp := rp.LocalProduct; rlp == nil {
	} else if len(rlp.Tags) > 0 {
		rlp.Tags = nil
	}
	return nil
}

func (pr *ProductRecord) AddTopNodes(topNodes models.Nodes, replace bool) error {
	return pr.LocalProduct.AddTopNodes(topNodes, replace)
}
