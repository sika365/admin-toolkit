package category

func (i *Package) Migrator() error {
	return i.db.DB.AutoMigrate(
		&Category{},
		&CategoryRecord{},
	)
}
