package product

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gorm.io/gorm"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/pkg/image"
	"github.com/sika365/admin-tools/utils"
	"github.com/sirupsen/logrus"
)

type Repo interface {
	Create(ctx *context.Context, db *gorm.DB, products LocalProducts) error
	CreateRecord(ctx *context.Context, db *gorm.DB, prodRecs ProductRecords) error
	Read(ctx *context.Context, db *gorm.DB, filters url.Values) (MapProducts, error)
	ReadByBarcode(ctx *context.Context, db *gorm.DB, rec *ProductRecord, filters url.Values) (*ProductRecord, error)
	ReadImagesWithoutProduct(ctx *context.Context, db *gorm.DB, filters url.Values) (mimages image.MapImages, err error)
	Update(ctx *context.Context, db *gorm.DB, product *LocalProduct, filters url.Values) error
	Delete(ctx *context.Context, db *gorm.DB, id simutils.PID, filters url.Values) error
}

type repo struct {
	client *client.Client
}

func newRepo(client *client.Client) (Repo, error) {
	r := &repo{
		client: client,
	}
	return r, nil
}

// Create implements Repo.
func (i *repo) Create(ctx *context.Context, db *gorm.DB, products LocalProducts) error {
	if len(products) == 0 {
		return nil
	} else if err := db.CreateInBatches(products, 100).Error; err != nil {
		return err
	} else {
		return nil
	}
}

// Create stores product records
func (i *repo) CreateRecord(ctx *context.Context, db *gorm.DB, prodRecs ProductRecords) error {
	if err := db.CreateInBatches(prodRecs, 100).Error; err != nil {
		return err
	} else {
		return nil
	}
}

// Read reads products with filters
func (i *repo) Read(ctx *context.Context, db *gorm.DB, filters url.Values) (products MapProducts, err error) {
	var stored LocalProducts
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return NewMapProducts(stored...), nil
	}
}

// ReadByBarcode reads products bye barcode
func (i *repo) ReadByBarcode(ctx *context.Context, db *gorm.DB, rec *ProductRecord, filters url.Values) (*ProductRecord, error) {
	var stored ProductRecord
	if err := db.
		Preload("LocalProduct").
		Preload("LocalProduct.Cover").
		Preload("LocalProduct.Gallery").
		Preload("LocalProduct.Product").
		Preload("LocalProduct.Product.Nodes").
		Preload("LocalCategory.Cover.Image").
		Preload("LocalCategory.Cover.File").
		Preload("LocalCategory.Category").
		Preload("LocalCategory.Category.Nodes").
		Preload("LocalCategory.Nodes").
		Where("barcode = ?", rec.Barcode).
		Take(&stored).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// register from remote
		if product, err := i.client.GetProductbyBarcode(ctx, rec.Barcode, filters); err != nil {
			return nil, err
		} else if rec.LocalProduct = FromProduct(product); rec.LocalProduct == nil {
			return nil, fmt.Errorf("nil local product")
		} else if err := i.CreateRecord(ctx, db, ProductRecords{rec}); err != nil {
			return nil, err
		} else {
			return rec, nil
		}
	} else if err != nil {
		return nil, err
	} else {
		return &stored, nil
	}
}

func (i *repo) ReadImagesWithoutProduct(ctx *context.Context, db *gorm.DB, filters url.Values) (mimages image.MapImages, err error) {
	var images image.LocalImages
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		InnerJoins("File").
		InnerJoins("Image").
		Joins("LEFT JOIN product_images ON local_images.id = product_images.local_image_id AND product_images.deleted_at IS NULL").
		Where("product_images.local_product_id IS NULL").
		// Where("images.title REGEXP '^[0-9]+$'").
		Find(&images).Error; err != nil {
		return nil, err
	} else if barcodeRegex, err := regexp.Compile(image.ImageBarcodeRegex); err != nil {
		return nil, err
	} else {
		var barcodeImages image.LocalImages
		for _, img := range images {
			if barcodeRegex.MatchString(img.Image.Title) {
				barcodeImages = append(barcodeImages, img)
			}
		}
		return image.NewMapImages(barcodeImages...), nil
	}
}

// Update implements Repo.
func (i *repo) Update(ctx *context.Context, db *gorm.DB, lprod *LocalProduct, filters url.Values) error {
	if rprod, err := i.client.PutProduct(ctx, lprod.Product); err != nil {
		return err
	} else if lprod.Product = rprod; false {
		return nil
	} else if err := db.
		Model(&LocalProduct{
			CommonTableFields: models.CommonTableFields{Model: database.Model{ID: lprod.ID}},
		}).
		Updates(lprod).Error; err != nil {
		return err
	} else if err := db.
		Model(lprod).
		Association("Gallery").
		Replace(lprod.Gallery); err != nil {
		return err
	} else if err := db.
		Model(lprod).
		Association("Product").
		Replace(lprod.Product); err != nil {
		return err
	} else {
		logrus.Infof("%s Updated", lprod.Product.LocalProduct.Barcodes)
		return nil
	}
}

// Delete implements Repo.
func (i *repo) Delete(ctx *context.Context, db *gorm.DB, id simutils.PID, filters url.Values) error {
	panic("unimplemented")
}
