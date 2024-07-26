package product

import (
	"net/url"
	"regexp"

	"gorm.io/gorm"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/image"
	"github.com/sika365/admin-tools/utils"
)

type Repo interface {
	Create(ctx *context.Context, db *gorm.DB, products LocalProducts) error
	CreateRecord(ctx *context.Context, db *gorm.DB, prodRecs ...*ProductRecord) error
	Read(ctx *context.Context, db *gorm.DB, filters url.Values) (MapProducts, error)
	ReadByBarcode(ctx *context.Context, db *gorm.DB, barcode string) (MapProducts, error)
	ReadImagesWithoutProduct(ctx *context.Context, db *gorm.DB, filters url.Values) (mimages image.MapImages, err error)
	Update(ctx *context.Context, db *gorm.DB, product *LocalProduct, filters url.Values) error
	Delete(ctx *context.Context, db *gorm.DB, id simutils.PID, filters url.Values) error
}

type repo struct {
}

func newRepo() (Repo, error) {
	r := &repo{}
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
func (i *repo) CreateRecord(ctx *context.Context, db *gorm.DB, prodRecs ...*ProductRecord) error {
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
func (i *repo) ReadByBarcode(ctx *context.Context, db *gorm.DB, barcode string) (MapProducts, error) {
	var stored LocalProducts
	if err := db.
		// Where("all_barcodes like ?", "'%"+barcode+";%'").
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return NewMapProducts(stored...), nil
	}
}

func (i *repo) ReadImagesWithoutProduct(ctx *context.Context, db *gorm.DB, filters url.Values) (mimages image.MapImages, err error) {
	var images image.LocalImages
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		InnerJoins("File").
		Joins("LEFT JOIN product_images ON images.id = product_images.image_id AND product_images.deleted_at IS NULL").
		Where("product_images.product_id IS NULL").
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
func (i *repo) Update(ctx *context.Context, db *gorm.DB, product *LocalProduct, filters url.Values) error {
	panic("unimplemented")
}

// Delete implements Repo.
func (i *repo) Delete(ctx *context.Context, db *gorm.DB, id simutils.PID, filters url.Values) error {
	panic("unimplemented")
}
