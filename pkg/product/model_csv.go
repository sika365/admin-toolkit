package product

import (
	"github.com/sika365/admin-tools/pkg/category"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
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
