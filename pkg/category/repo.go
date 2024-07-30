package category

import (
	"net/url"

	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gorm.io/gorm"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/utils"
)

type Repo interface {
	Create(ctx *context.Context, db *gorm.DB, value any) error
	Read(ctx *context.Context, db *gorm.DB, filters url.Values) (LocalCategories, error)
	Update(ctx *context.Context, db *gorm.DB, category *LocalCategory, filters url.Values) error
	Delete(ctx *context.Context, db *gorm.DB, id database.PID, filters url.Values) error
	Clear(ctx *context.Context, db *gorm.DB) error
}

type repo struct {
}

func newRepo() (Repo, error) {
	r := &repo{}
	return r, nil
}

// Create ...
func (i *repo) Create(ctx *context.Context, db *gorm.DB, value any) error {
	if err := db.CreateInBatches(value, 100).Error; err != nil {
		return err
	} else {
		return nil
	}
}

// Read fetches categories with filters
func (i *repo) Read(ctx *context.Context, db *gorm.DB, filters url.Values) (LocalCategories, error) {
	var stored LocalCategories
	if err := utils.
		BuildGormQuery(ctx, db, filters).
		Preload("Cover.Image").
		Preload("Cover.File").
		Preload("Category").
		Preload("Category.Nodes").
		Preload("Nodes").
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return stored, nil
	}
}

// Update implements Repo.
func (i *repo) Update(ctx *context.Context, db *gorm.DB, category *LocalCategory, filters url.Values) error {
	panic("unimplemented")
}

// Delete implements Repo.
func (i *repo) Delete(ctx *context.Context, db *gorm.DB, id database.PID, filters url.Values) error {
	panic("unimplemented")
}

func (i *repo) Clear(ctx *context.Context, db *gorm.DB) error {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&CategoryRecord{}).Error; err != nil {
		return err
	}
	return nil
}
