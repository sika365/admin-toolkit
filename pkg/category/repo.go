package category

import (
	"net/url"

	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gorm.io/gorm"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/utils"
)

type Repo interface {
	Create(ctx *context.Context, db *gorm.DB, catRecs CategoryRecords) error
	ReadCategoryRecords(ctx *context.Context, db *gorm.DB, filters url.Values) (categories CategoryRecords, err error)
	Read(ctx *context.Context, db *gorm.DB, filters url.Values) (LocalCategories, error)
	ReadCategory(ctx *context.Context, db *gorm.DB, filters url.Values) (categories models.Categories, err error)
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

// Create creates category records in batches
func (i *repo) Create(ctx *context.Context, db *gorm.DB, catRecs CategoryRecords) error {
	for _, catRec := range catRecs {
		if err := db.FirstOrCreate(catRec).Error; err != nil {
			return err
		}
	}

	return nil
}

// Read reads categories with filters
func (i *repo) ReadCategoryRecords(ctx *context.Context, db *gorm.DB, filters url.Values) (categories CategoryRecords, err error) {
	var stored CategoryRecords
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		Preload("LocalCategory").
		Preload("LocalCategory.Cover").
		Preload("LocalCategory.Cover.Image").
		Preload("LocalCategory.Cover.File").
		Preload("LocalCategory.Nodes").
		Preload("LocalCategory.Nodes.Parent").
		Preload("LocalCategory.Nodes.SubNodes").
		Preload("LocalCategory.Category").
		Preload("LocalCategory.Category.Nodes").
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return stored, nil
	}
}

// Read fetches categories with filters
func (i *repo) Read(ctx *context.Context, db *gorm.DB, filters url.Values) (LocalCategories, error) {
	var stored LocalCategories
	if err := utils.
		BuildGormQuery(ctx, db, filters).
		Preload("Cover").
		Preload("Cover.Image").
		Preload("Cover.File").
		Preload("Nodes").
		Preload("Nodes.Parent").
		Preload("Nodes.SubNodes").
		Preload("Category").
		Preload("Category.Nodes").
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return stored, nil
	}
}

func (i *repo) ReadCategory(ctx *context.Context, db *gorm.DB, filters url.Values) (categories models.Categories, err error) {
	if err := utils.
		BuildGormQuery(ctx, db, filters).
		Preload("Nodes").
		Preload("Nodes.Parent").
		Find(&categories).Error; err != nil {
		return nil, err
	} else {
		return categories, nil
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
