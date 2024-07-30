package product

import (
	"net/url"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alifakhimi/simple-utils-go/simscheme"
	"github.com/alitto/pond"
	"github.com/sirupsen/logrus"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/config"
	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/category"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/pkg/excel"
	"github.com/sika365/admin-tools/pkg/image"
)

type Logic interface {
	Find(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (products LocalProducts, err error)
	Create(ctx *context.Context, prods LocalProducts, batchSize int) error
	SyncByImages(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (*simscheme.Document, error)
	SyncBySpreadSheets(ctx *context.Context, req *SyncBySpreadSheetsRequest, filters url.Values) (*simscheme.Document, error)
	SetImage(ctx *context.Context, req *SyncByImageRequest, rec *ProductRecord, limages image.LocalImages) error
	MatchBarcode(ctx *context.Context, req *SyncByImageRequest, barcode string, productNodes models.Nodes) (*models.Product, error)
}

type logic struct {
	conn     *simutils.DBConnection
	client   *client.Client
	repo     Repo
	catLogic category.Logic
	catRepo  category.Repo
}

func newLogic(conn *simutils.DBConnection, client *client.Client, repo Repo, catLogic category.Logic, catRepo category.Repo) (Logic, error) {
	l := &logic{
		conn:     conn,
		client:   client,
		repo:     repo,
		catLogic: catLogic,
		catRepo:  catRepo,
	}
	return l, nil
}

func (l *logic) Find(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (products LocalProducts, err error) {
	q := l.conn.DB.WithContext(ctx.Request().Context())

	if products, err := l.repo.Read(ctx, q, filters); err != nil {
		return nil, err
	} else {
		return products.GetValues(), nil
	}
}

func (l *logic) Create(ctx *context.Context, prods LocalProducts, batchSize int) error {
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
			} else if err := l.repo.Create(ctx, tx, LocalProducts{&LocalProduct{Product: resultProd}}); err != nil {
				logrus.Infof("writing file %v in db failed", prd)
				return
			}
		})
	}

	pool.StopAndWait()

	return nil
}

func (l *logic) SyncByImages(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (*simscheme.Document, error) {
	var (
		err              error
		batchSize        = 1
		mimages          image.MapImages
		mapBarcodeImages = make(map[string]image.LocalImages)
		pool             = pond.New(batchSize, 0)
		prodRecDoc       = simscheme.
					GetSchema().
					AddNewDocumentWithType(&ProductRecord{})
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
				prodRec = &ProductRecord{Barcode: barcode}
			)
			// https://sika365.com/admin/api/v1/nodes/root/products?order_by=newest&search=7899665999353&check_availability=false&search_products_in_nodes=true&search_in_node=false&search_in_sub_node=false&get_product_parents=false&search_in_reserved_quantity=false&search_in_limited_quantity=false&coverstatus=0&total=0&limit=20&offset=0&cover_status=-1&view=node&remote_pagination=false&remote_search=false&includes=Cover&includes=Nodes.Parent.Category&includes=Tags.Node.Category&includes=CategoryNodes&store_id=38&branch_id=47&stock_id=45

			// Is the image cover or for gallery?
			// Retrieve product by barcode
			if prd, err := l.repo.ReadByBarcode(ctx,
				l.conn.DB.WithContext(ctx.Request().Context()),
				prodRec,
				filters,
			); err != nil {
				return
			} else if err := l.SetImage(ctx, req, prd, imgs); err != nil {
				return
			} else if err := l.repo.Update(ctx,
				l.conn.DB.WithContext(ctx.Request().Context()),
				prd.LocalProduct,
				nil); err != nil {
				return
			} else {
				logrus.Infof("%s Updated", prd.Barcode)
				prodRecDoc.AddNode(prodRec)
			}
		})
	}

	pool.StopAndWait()

	return prodRecDoc, nil
}

func (l *logic) SyncBySpreadSheets(ctx *context.Context, req *SyncBySpreadSheetsRequest, filters url.Values) (*simscheme.Document, error) {
	var (
		prodRecDoc = simscheme.
				GetSchema().
				AddNewDocumentWithType(&ProductRecord{})

		prodNodeIDs database.PIDs
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
			var (
				err     error
				prodRec = &ProductRecord{
					Barcode:       rec[header[req.ProductHeaderMap.Barcode]],
					Title:         rec[header[req.ProductHeaderMap.Title]],
					CategoryAlias: rec[header[req.ProductHeaderMap.CategoryAlias]],
				}
			)

			if req.ProductHeaderMap.CategoryAlias != "" {
				if lcats, err := l.catRepo.Read(ctx,
					l.conn.DB.WithContext(ctx.Request().Context()),
					url.Values{
						"alias": []string{prodRec.CategoryAlias},
					},
				); err != nil {
					return
				} else if len(lcats) == 1 {
					prodRec.LocalCategory = lcats[0]
					for _, n := range prodRec.LocalCategory.Category.Nodes {
						prodNodeIDs = append(prodNodeIDs, n.ID)
					}
				}
			}

			if prodRec, err = l.repo.ReadByBarcode(ctx,
				l.conn.DB.WithContext(ctx.Request().Context()),
				prodRec,
				filters,
			); err != nil {
				return
			} else if len(prodNodeIDs) == 0 {
				// break
			} else if prodNodes, err := l.client.AddToNodes(ctx, prodRec.LocalProduct.Product, prodNodeIDs); err != nil {
				return
			} else if len(prodNodes) == 0 {
				// break
			} else {
				prodRec.LocalProduct.Product = prodNodes[0].Product
			}

			prodRecDoc.AddNode(prodRec)
		},
	); err != nil {
		return nil, err
	} else {
		return prodRecDoc, nil
	}
}

func (l *logic) SetImage(_ *context.Context, req *SyncByImageRequest, rec *ProductRecord, limages image.LocalImages) error {
	lprod := rec.LocalProduct
	rprod := lprod.Product

	if req.ReplaceGallery {
		clear(rprod.Images)
		rprod.Images = nil
	}

	for _, limg := range limages {
		// TODO if image is cover then set cover is true else add to gallery
		if (!lprod.CoverID.IsValid() && lprod.Cover == nil) &&
			(req.ReplaceCover || (!req.IgnoreCoverIfEmpty && !rprod.CoverID.IsValid())) {
			// product.Cover = &ProductImage{
			// 	ImageID:   img.ID,
			// 	Image:     img,
			// 	ProductID: product.ID,
			// }
			// product.CoverID, _ = img.ID.ToNullPID()
			lprod.Cover = limg
			rprod.CoverID, _ = limg.ID.ToNullPID()
		} else if req.ReplaceGallery || !req.IgnoreAddToGallery {
			lprod.Gallery = append(lprod.Gallery, &ProductImage{
				LocalImageID:   limg.ID,
				LocalImage:     limg,
				LocalProductID: lprod.ID,
			})
			rprod.Images = append(rprod.Images, &models.Imagable{
				ImageID: limg.ID,
				Image:   limg.Image,
			})
		}
	}

	return nil
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
