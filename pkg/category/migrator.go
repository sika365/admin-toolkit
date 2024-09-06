package category

import (
	"github.com/sika365/admin-tools/pkg/image"
	"github.com/sika365/admin-tools/pkg/node"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
)

func (i *Package) Migrator() error {
	return i.db.DB.AutoMigrate(
		&models.Node{},
		&models.Category{},
		&models.Imagable{},
		&image.LocalImage{},
		&node.LocalNode{},
		&LocalCategory{},
		&CategoryRecord{},
	)
}
