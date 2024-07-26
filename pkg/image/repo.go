package image

import (
	"net/url"

	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gorm.io/gorm"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/utils"
)

type Repo interface {
	Create(ctx *context.Context, db *gorm.DB, images LocalImages) error
	Read(ctx *context.Context, db *gorm.DB, filters url.Values) (MapImages, error)
	ReadFiles(ctx *context.Context, db *gorm.DB, files MapImages, filters url.Values) (MapImages, error)
	Update(ctx *context.Context, db *gorm.DB, file *LocalImage, filters url.Values) error
	Delete(ctx *context.Context, db *gorm.DB, id database.PID, filters url.Values) error
}

type repo struct {
}

func newRepo() (Repo, error) {
	r := &repo{}
	return r, nil
}

// Read reads files with filters
func (i *repo) Read(ctx *context.Context, db *gorm.DB, filters url.Values) (images MapImages, err error) {
	var stored LocalImages
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		InnerJoins("Image").
		InnerJoins("File").
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return NewMapImages(stored...), nil
	}
}

func (i *repo) ReadFiles(ctx *context.Context, db *gorm.DB, files MapImages, filters url.Values) (images MapImages, err error) {
	var imgs LocalImages
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		InnerJoins("Image").
		InnerJoins("File", db.Where("hash in (?)", files.GetKeys())).
		Find(&imgs).Error; err != nil {
		return nil, err
	} else {
		for _, img := range imgs {
			files[img.Hash()] = img
		}
		return files, nil
	}
}

// Create implements Repo.
func (i *repo) Create(ctx *context.Context, db *gorm.DB, files LocalImages) error {
	if len(files) == 0 {
		return nil
	} else if err := db.CreateInBatches(files, 100).Error; err != nil {
		return err
	} else {
		return nil
	}
}

// Update implements Repo.
func (i *repo) Update(ctx *context.Context, db *gorm.DB, file *LocalImage, filters url.Values) error {
	panic("unimplemented")
}

// Delete implements Repo.
func (i *repo) Delete(ctx *context.Context, db *gorm.DB, id database.PID, filters url.Values) error {
	panic("unimplemented")
}
