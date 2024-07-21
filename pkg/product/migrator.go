package product

func (i *impl) Migrator() error {
	return i.db.DB.AutoMigrate(
		&ProductImage{},
		&Product{},
	)
}
