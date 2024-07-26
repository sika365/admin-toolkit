package product

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"

	"github.com/sika365/admin-tools/pkg/category"
)

type ProductRecord struct {
	Barcode        string                   `json:"barcode,omitempty" gorm:"primaryKey"`
	Title          string                   `json:"title,omitempty"`
	Category       string                   `json:"category,omitempty"`
	CategoryRecord *category.CategoryRecord `json:"category_record,omitempty" gorm:"foreignKey:Category;references:Title"`
	ProductID      database.PID             `json:"product_id,omitempty"`
	Product        *LocalProduct            `json:"product,omitempty"`
}
