package file

func (i *impl) Migrator() error {
	return i.db.DB.AutoMigrate(&File{})
}
