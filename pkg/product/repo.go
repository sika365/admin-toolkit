package product

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/pkg/image"
	"github.com/sika365/admin-tools/utils"
	"github.com/sirupsen/logrus"
)

type Repo interface {
	Save(ctx *context.Context, db *gorm.DB, products LocalProducts) error
	Create(ctx *context.Context, db *gorm.DB, products LocalProducts) error
	CreateProductGroup(ctx *context.Context, db *gorm.DB, prdgrp *LocalProductGroup) error
	CreateRecord(ctx *context.Context, db *gorm.DB, prodRecs ProductRecords) error
	Read(ctx *context.Context, db *gorm.DB, filters url.Values) (MapProducts, error)
	FirstOrCreateLocalProuctGroup(ctx *context.Context, db *gorm.DB, reqLProductGroup *LocalProductGroup) (*LocalProductGroup, error)
	ReadProductRecords(ctx *context.Context, db *gorm.DB, filters url.Values) (products ProductRecords, err error)
	ReadByBarcode(ctx *context.Context, db *gorm.DB, rec *ProductRecord, filters url.Values) (*ProductRecord, error)
	ReadImagesWithoutProduct(ctx *context.Context, db *gorm.DB, filters url.Values) (mimages image.MapImages, err error)
	UpdateImages(ctx *context.Context, db *gorm.DB, product *LocalProduct, filters url.Values) error
	UpdateNodes(ctx *context.Context, db *gorm.DB, prdRec *ProductRecord, topNodes models.Nodes, filters url.Values) error
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

func (i *repo) CreateProductGroup(ctx *context.Context, db *gorm.DB, lprdgrp *LocalProductGroup) error {
	if lprdgrp == nil {
		return nil
	} else if err := db.Create(lprdgrp).Error; err != nil {
		return err
	} else {
		return nil
	}
}

func (i *repo) Save(ctx *context.Context, db *gorm.DB, products LocalProducts) error {
	for _, prd := range products {
		if err := db.Save(prd).Error; err != nil {
			return err
		} else {
			return nil
		}
	}
	return nil
}

// Create ...
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
	for _, prodRec := range prodRecs {
		if prodRec.LocalProduct != nil && prodRec.LocalProduct.Product != nil &&
			(prodRec.LocalProduct.Product.LocalProduct == nil ||
				!database.IsValid(prodRec.LocalProduct.Product.ProductStock.ProductID)) {
			// Register product
			if prd, err := i.client.CreateProduct(
				ctx,
				prodRec.LocalProduct.Product,
			); err != nil {
				return err
			} else {
				prodRec.LocalProduct.Product = prd
			}
		}
	}

	if len(prodRecs) == 0 {
		return nil
	} else if err := db.CreateInBatches(prodRecs, 100).Error; err != nil {
		return err
	} else {
		return nil
	}
}

// Read reads products with filters
func (i *repo) ReadProductRecords(ctx *context.Context, db *gorm.DB, filters url.Values) (products ProductRecords, err error) {
	var stored ProductRecords
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return stored, nil
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

func (i *repo) FirstOrCreateLocalProuctGroup(ctx *context.Context, db *gorm.DB, reqLProductGroup *LocalProductGroup) (*LocalProductGroup, error) {
	var (
		storedPrdGrp LocalProductGroup
		// response     = models.ProductGroupResponse{}
		// request      = models.ProductGroupRequest{ProductGroup: *prdgrp.ProductGroup}
	)
	// 1) First or create the product group
	// 1-1) Fetch from db
	if err := db.
		Joins("ProductGroup").
		// Preload("ProductGroup", "slug=?", lprdgrp.ProductGroup.Slug).
		// Joins("ProductGroup.Cover").Preload("ProductGroup.Images").
		Preload("ProductGroup.Products").Preload("ProductGroup.Products.Nodes").
		// Preload("ProductGroup.Products.Cover").Preload("ProductGroup.Products.Images").
		Joins("Cover").Preload("Gallery").Preload("Gallery.Image").
		Where(&LocalProductGroup{Slug: reqLProductGroup.ProductGroup.Slug}).
		Take(&storedPrdGrp).Error; errors.Is(err, gorm.ErrRecordNotFound) ||
		storedPrdGrp.ProductGroup == nil {
		var rprdgrp *models.ProductGroup
		// Retrieve from remote
		if rprdgrp, err = i.client.GetProductGroupBySlug(
			ctx,
			simutils.MakeSlug(reqLProductGroup.ProductGroup.Slug),
		); err != nil && !errors.Is(err, models.ErrNotFound) {
			return nil, err
		} else if errors.Is(err, models.ErrNotFound) {
			// register
			if rprdgrp, err = i.client.CreateProductGroup(
				ctx,
				reqLProductGroup.ProductGroup,
			); err != nil {
				return nil, err
			}
		}

		if rprdgrp == nil {
			return nil, models.ErrNotFound
		}

		reqProductGroupProducts := reqLProductGroup.ProductGroup.Products
		reqLProductGroup.ProductGroup = rprdgrp
		reqLProductGroup.ProductGroup.Products = reqProductGroupProducts

		// Write product group into the database
		if tx := db.WithContext(ctx.Request().Context()); tx == nil {
			return nil, err
		} else if err := i.CreateProductGroup(ctx, tx, reqLProductGroup); err != nil {
			logrus.Infof("create product group in db failed %v", err)
			return nil, err
		} else {
			return reqLProductGroup, nil
		}
	} else {
		return &storedPrdGrp, nil
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
		// Preload("LocalCategory.Nodes").
		Where("barcode = ?", rec.Barcode).
		Take(&stored).Error; errors.Is(err, gorm.ErrRecordNotFound) ||
		stored.LocalProduct == nil ||
		stored.LocalProduct.Product == nil {
		// register from remote
		if products, err := i.client.GetProductsByBarcode(ctx, rec.Barcode, filters); err != nil {
			return nil, err
		} else if product, err := models.MergeProductStocks(products); err != nil {
			return nil, err
		} else if err = func(rec *ProductRecord, p *models.Product) error {
			if rec.LocalProduct == nil || rec.LocalProduct.Product == nil {
				// logrus.WithFields(logrus.Fields{
				// 	"barcode":        rec.Barcode,
				// 	"product_record": rec,
				// }).Errorln(ErrRemoteProductNotFound)
				// return ErrRemoteProductNotFound
				rec.LocalProduct = FromProduct(p)
				return nil
			} else {
				rprod := rec.LocalProduct.Product
				rprod.ID = p.ID
				rprod.LocalProduct.ID = p.LocalProduct.ID
				rprod.LocalProduct.StoreID = p.ProductStock.StoreID
				rprod.ProductStock = p.ProductStock
				rprod.LocalProduct.ProductStocks = p.ProductStocks
				return nil
			}
		}(rec, product); err != nil {
			return nil, err
			// } else if rec.LocalProduct = FromProduct(product); rec.LocalProduct == nil {
			// 	return nil, fmt.Errorf("nil local product")
		} else if err := i.CreateRecord(ctx, db, ProductRecords{rec}); err != nil {
			return nil, err
		} else {
			return rec, nil
		}
	} else if err != nil {
		logrus.WithFields(logrus.Fields{
			"fn":             "Product.repo.ReadByBarcode",
			"product_record": rec,
		}).Errorln(err)
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
		Where("local_images.image_id not in (?)",
			db.Table("viw_products").Select("viw_products.cover_id").
				Where("viw_products.cover_id=local_images.image_id"),
		).
		Order("local_images.id asc").
		// Joins("LEFT JOIN products ON products.cover_id = local_images.image_id AND products.deleted_at IS NULL").
		// Where("products.cover_id IS NULL").
		// Joins("LEFT JOIN images ON local_images.image_id = images.id AND images.deleted_at IS NULL").
		// Joins("LEFT JOIN products ON products.cover_id = images.id AND products.deleted_at IS NULL").
		// Joins("LEFT JOIN local_products ON local_products.cover_id = local_images.id AND local_products.deleted_at IS NULL").
		// Joins("LEFT JOIN product_images ON local_images.id = product_images.local_image_id AND product_images.deleted_at IS NULL").
		// Where("product_images.local_product_id IS NULL").
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

func (i *repo) UpdateImages(ctx *context.Context, db *gorm.DB, lprod *LocalProduct, filters url.Values) error {
	if lprod.Product == nil {
		logrus.Warnf("local product id %s remote product not found", lprod.ID)
		return nil
	} else if !lprod.Product.CoverID.IsValid() {
		logrus.Warnf("local product id %s cover not found", lprod.ID)
		return nil
	} else if err := i.Update(ctx, db, lprod, filters); err != nil {
		return err
	} else if err := db.
		Model(lprod).
		Association("Gallery").
		Replace(lprod.Gallery); err != nil {
		return err
	} else {
		return nil
	}
}

func (i *repo) UpdateNodes(ctx *context.Context, db *gorm.DB, prdRec *ProductRecord, topNodes models.Nodes, filters url.Values) error {
	if req, err := context.GetRequestModel[*SyncBySpreadSheetsRequest](ctx); err != nil {
		return err
	} else if prdRec == nil {
		return fmt.Errorf("product record is nil")
	} else if lprod := prdRec.LocalProduct; lprod == nil || prdRec.LocalProduct.Product == nil {
		logrus.Warnf("local product id %s remote product not found", prdRec.Barcode)
		return nil
	} else if len(topNodes) == 0 {
		logrus.Warnf("no top nodes found for %s", prdRec.Barcode)
		return nil
	} else if err := lprod.AddTopNodes(topNodes, req.ReplaceNodes); err != nil {
		logrus.Warnf("add top nodes error for %s => %v", prdRec.Barcode, err)
		return nil
	} else if err := i.Update(ctx, db, lprod, filters); err != nil {
		return err
	} else {
		return nil
	}
}

// Update implements Repo.
func (i *repo) Update(ctx *context.Context, db *gorm.DB, lprod *LocalProduct, filters url.Values) (err error) {
	if lprod.Product == nil {
		logrus.Warnf("local product id %s remote product not found", lprod.ID)
		return nil
	} else if lprod.Product, err = i.client.PutProduct(ctx, lprod.Product); err != nil {
		return err
	} else if err := db.
		Model(&LocalProduct{
			CommonTableFields: models.CommonTableFields{Model: database.Model{ID: lprod.ID}},
		}).
		Updates(lprod).Error; err != nil {
		return err
	} else if err := db.Unscoped().
		// Model(lprod).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		// Association("Product").
		Omit(clause.Associations).
		Save(lprod.Product).Error; err != nil {
		return err
	} else if err := db.Unscoped().
		Model(lprod.Product).
		// Session(&gorm.Session{FullSaveAssociations: true}).
		Association("Nodes").
		Unscoped().Replace(lprod.Product.Nodes); err != nil {
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
