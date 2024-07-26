package file

func (i *Package) Migrator() error {
	return i.db.DB.AutoMigrate(&File{})
}
