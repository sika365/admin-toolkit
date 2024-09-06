package woocommerce

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/pkg/category"
	"github.com/sika365/admin-tools/pkg/image"
	"github.com/sika365/admin-tools/pkg/node"
)

func (i *Package) Migrator() error {
	return i.db.DB.AutoMigrate(
		&models.Node{},
		&models.Category{},
		&models.ProductGroup{},
		&image.LocalImage{},
		&node.LocalNode{},
		&category.LocalCategory{},
		&category.CategoryRecord{},
	)
}
