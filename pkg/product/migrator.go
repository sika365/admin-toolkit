package product

func (i *Package) Migrator() error {
	return i.db.DB.AutoMigrate(
		&ProductImage{},
		&Product{},
	)
}
