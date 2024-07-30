package node

func (i *Package) Migrator() error {
	return i.db.DB.AutoMigrate(
		&LocalNode{},
	)
}
