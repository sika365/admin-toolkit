package product

import (
	"net/url"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alifakhimi/simple-utils-go/simscheme"
	"github.com/alitto/pond"
	"github.com/sirupsen/logrus"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/config"
	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/pkg/excel"
	"github.com/sika365/admin-tools/pkg/image"
)

type Logic interface {
	Find(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (products Products, err error)
	Create(ctx *context.Context, prods Products, batchSize int) error
	SyncByImages(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (Products, error)
	SyncBySpreadSheets(ctx *context.Context, req *SyncBySpreadSheetsRequest, filters url.Values) (*simscheme.Document, error)
	SetImage(ctx *context.Context, req *SyncByImageRequest, prd *models.Product, imgs image.Images) (*Product, error)
	MatchBarcode(ctx *context.Context, req *SyncByImageRequest, barcode string, productNodes models.Nodes) (*models.Product, error)
}

type logic struct {
	conn   *simutils.DBConnection
	client *client.Client
	repo   Repo
}

func newLogic(conn *simutils.DBConnection, client *client.Client, repo Repo) (Logic, error) {
	l := &logic{
		conn:   conn,
		client: client,
		repo:   repo,
	}
	return l, nil
}

func (l *logic) Find(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (products Products, err error) {
	q := l.conn.DB.WithContext(ctx.Request().Context())

	if products, err := l.repo.Read(ctx, q, filters); err != nil {
		return nil, err
	} else {
		return products.GetValues(), nil
	}
}

func (l *logic) Create(ctx *context.Context, prods Products, batchSize int) error {
	pool := pond.New(batchSize, 0)

	for _, prd := range prods {
		pool.Submit(func() {
			var (
				conf         = config.Config()
				productsResp = models.ProductsResponse{}
			)

			logrus.Infof("Running task for %v", prd)
			// Upload files
			if client, err := conf.GetRestyClient("sika365"); err != nil {
				return
			} else if resp, err := client.R().
				SetBody(prd).
				SetResult(&productsResp).
				SetError(&productsResp).
				Put("/products"); err != nil {
				logrus.Info(err)
				return
			} else if !resp.IsSuccess() {
				return
			} else if prods := productsResp.Data.Products; len(prods) == 0 || prods[0] == nil {
				return
			} else if resultProd := prods[0]; false {
				return
				// Write uploaded files into the database
			} else if tx := l.conn.DB.WithContext(ctx.Request().Context()); tx == nil {
				return
			} else if err := l.repo.Create(ctx, tx, Products{&Product{Product: resultProd}}); err != nil {
				logrus.Infof("writing file %v in db failed", prd)
				return
			}
		})
	}

	pool.StopAndWait()

	return nil
}

func (l *logic) SyncByImages(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (products Products, err error) {
	var (
		batchSize        = 5
		filtersEncoded   = filters.Encode()
		mimages          image.MapImages
		mapBarcodeImages = make(map[string]image.Images)
		pool             = pond.New(batchSize, 0)
	)

	// Get barcodes from image's title if synced and there aren't any related product
	if mimages, err = l.repo.ReadImagesWithoutProduct(
		ctx,
		l.conn.DB.WithContext(ctx.Request().Context()),
		url.Values{},
	); err != nil {
		return nil, err
	}

	for _, img := range mimages {
		// TODO check barcode pattern
		mapBarcodeImages[img.Image.Title] = append(mapBarcodeImages[img.Image.Title], img)
	}

	for barcode, imgs := range mapBarcodeImages {
		pool.Submit(func() {
			var (
				filters, _ = url.ParseQuery(filtersEncoded)
				conf       = config.Config()
				code       = barcode
				// tx           = l.conn.DB.WithContext(ctx.Request().Context())
				productsResp      = ProductSearchResponse{}
				updateProductResp = models.ProductsResponse{}
			)
			// https://sika365.com/admin/api/v1/nodes/root/products?order_by=newest&search=7899665999353&check_availability=false&search_products_in_nodes=true&search_in_node=false&search_in_sub_node=false&get_product_parents=false&search_in_reserved_quantity=false&search_in_limited_quantity=false&coverstatus=0&total=0&limit=20&offset=0&cover_status=-1&view=node&remote_pagination=false&remote_search=false&includes=Cover&includes=Nodes.Parent.Category&includes=Tags.Node.Category&includes=CategoryNodes&store_id=38&branch_id=47&stock_id=45
			filters.Set("search", code)

			// Is the image cover or for gallery?
			// Retrieve product by barcode
			if client, err := conf.GetRestyClient("sika365"); err != nil {
				return
			} else if resp, err := client.R().
				SetQueryParamsFromValues(filters).
				SetResult(&productsResp).
				SetError(&productsResp).
				Get("/nodes/root/products"); err != nil {
				logrus.Info(err)
				return
			} else if !resp.IsSuccess() {
				return
			} else if prd, err := l.MatchBarcode(ctx, req, code, productsResp.Data.ProductNodes); err != nil {
				return
			} else if prd, err := l.SetImage(ctx, req, prd, imgs); err != nil {
				return
			} else if resp, err := client.R().
				SetPathParams(map[string]string{
					"id": prd.Product.ID.String(),
				}).
				SetBody(prd.Product).
				SetResult(&updateProductResp).
				SetError(&updateProductResp).
				Put("/products/{id}"); err != nil {
				return
			} else if !resp.IsSuccess() {
				return
			} else if err := l.repo.Create(
				ctx,
				l.conn.DB.WithContext(ctx.Request().Context()),
				Products{prd},
			); err != nil {
				return
			} else {
				logrus.Infof("%s Updated", prd.Product.LocalProduct.Barcodes)
				products = append(products, prd)
			}
		})
	}

	pool.StopAndWait()

	return products, nil
}

func (l *logic) SetImage(_ *context.Context, req *SyncByImageRequest, prd *models.Product, imgs image.Images) (*Product, error) {
	product := FromProduct(prd)

	if req.ReplaceGallery {
		clear(prd.Images)
		prd.Images = nil
	}

	for _, img := range imgs {
		// TODO if image is cover then set cover is true else add to gallery
		if (!product.CoverID.IsValid() && product.Cover == nil) &&
			(req.ReplaceCover || (!req.IgnoreCoverIfEmpty && !prd.CoverID.IsValid())) {
			// product.Cover = &ProductImage{
			// 	ImageID:   img.ID,
			// 	Image:     img,
			// 	ProductID: product.ID,
			// }
			// product.CoverID, _ = img.ID.ToNullPID()
			product.Cover = img
			prd.CoverID, _ = img.ID.ToNullPID()
		} else if req.ReplaceGallery || !req.IgnoreAddToGallery {
			product.Gallery = append(product.Gallery, &ProductImage{
				ImageID:   img.ID,
				Image:     img,
				ProductID: product.ID,
			})
			prd.Images = append(prd.Images, &models.Imagable{
				ImageID: img.ID,
				// Image:     pi.Image.Image,
			})
		}
	}

	return product, nil
}

func (l *logic) MatchBarcode(_ *context.Context, _ *SyncByImageRequest, barcode string, productNodes models.Nodes) (*models.Product, error) {
	for _, node := range productNodes {
		if node.Product == nil || node.Product.LocalProduct == nil {
			continue
		}
		for _, b := range node.Product.LocalProduct.Barcodes {
			if b.Barcode == barcode {
				return node.Product, nil
			}
		}
	}

	return nil, simutils.ErrNotFound
}

func (l *logic) SyncBySpreadSheets(ctx *context.Context, req *SyncBySpreadSheetsRequest, filters url.Values) (*simscheme.Document, error) {
	var (
		prodRecDoc = simscheme.
				GetSchema().
				AddNewDocumentWithType(&ProductRecord{})

		productRecords []*ProductRecord
	)

	if req.ProductHeaderMap.Barcode == "" {
		return nil, nil
	} else if csvFiles, err := excel.LoadExcels(ctx, req.Root, req.MaxDepth); err != nil {
		return nil, err
		// Make ProductNodes from the files
	} else if err := excel.FromFiles(
		csvFiles,
		req.Offset,
		func(header map[string]int, rec []string) {
			prodRec := &ProductRecord{
				Barcode:        rec[header[req.ProductHeaderMap.Barcode]],
				Title:          rec[header[req.ProductHeaderMap.Title]],
				Category:       rec[header[req.ProductHeaderMap.Category]],
				CategoryRecord: nil,
			}
			prodRecDoc.AddNode(prodRec)
		},
	); err != nil {
		return prodRecDoc, err
	} else if err := prodRecDoc.GetData(prodRecDoc.GetData(productRecords)); err != nil {
		return prodRecDoc, err
	} else if err := l.repo.CreateRecord(ctx,
		l.conn.DB.WithContext(ctx.Request().Context()),
		productRecords...,
	); err != nil {
		return nil, err
	} else {
		return prodRecDoc, nil
	}
}
