package product

import "gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

func (i *Package) Migrator() error {
	return i.db.DB.AutoMigrate(
		&models.UnitType{},
		&models.Product{},
		&ProductImage{},
		&LocalProduct{},
		&ProductRecord{},
	)
}
