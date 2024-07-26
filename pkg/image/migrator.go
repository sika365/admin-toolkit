package image

import (
	"github.com/sika365/admin-tools/pkg/file"
	"github.com/sika365/admin-tools/registrar"
)

func (i *Package) Migrator() error {
	if fp, err := registrar.Get(file.PackageName); err != nil {
		return err
	} else if err := fp.Migrator(); err != nil {
		return err
	} else {
		return i.db.DB.AutoMigrate(&Image{})
	}
}
