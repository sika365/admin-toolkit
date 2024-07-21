package file

func (i *pkg) Migrator() error {
	return i.db.DB.AutoMigrate(&File{})
}
