package product

import (
	"fmt"
	"net/url"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alifakhimi/simple-utils-go/simscheme"
	"github.com/alitto/pond"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/category"
	"github.com/sika365/admin-tools/pkg/excel"
	"github.com/sika365/admin-tools/pkg/image"
	"github.com/sika365/admin-tools/service/client"
)

type Logic interface {
	Find(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (products LocalProducts, err error)
	Save(ctx *context.Context, reqPrdRec *ProductRecord) (prdRec *ProductRecord, err error)
	FindOrCreateProductGroup(ctx *context.Context, lprdgrp *LocalProductGroup) (*LocalProductGroup, error)
	SyncByImages(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (*simscheme.Document, error)
	SyncBySpreadSheets(ctx *context.Context) (*simscheme.Document, error)
	SetImages(ctx *context.Context, req *SyncByImageRequest, rec *ProductRecord, limages image.LocalImages) error
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

func (l *logic) UpdateNodes(ctx *context.Context, prodRec *ProductRecord) (err error) {
	var (
		topNodes models.Nodes
	)

	if prodRec.CategorySlug != "" {
		if lcats, err := l.catRepo.Read(ctx,
			l.conn.DB.WithContext(ctx.Request().Context()),
			url.Values{
				"slug": []string{prodRec.CategorySlug},
			},
		); err != nil {
			return err
		} else if len(lcats) == 1 {
			prodRec.LocalCategory = lcats[0]
			topNodes = append(topNodes, prodRec.LocalCategory.Category.Nodes...)
		}
	} else {
		return nil
	}

	if prodRec, err = l.repo.ReadByBarcode(ctx,
		l.conn.DB.WithContext(ctx.Request().Context()),
		prodRec,
		ctx.QueryParams(),
	); err != nil {
		return err
	} else if err := l.repo.UpdateNodes(ctx,
		l.conn.DB.WithContext(ctx.Request().Context()),
		prodRec,
		topNodes,
		ctx.QueryParams(),
	); err != nil {
		return err
	} else {
		return nil
	}
}

func (l *logic) UpdateImage(ctx *context.Context, req *SyncByImageRequest, imgs image.LocalImages, prodRec *ProductRecord) (err error) {
	// Is the image cover or for gallery?
	// Retrieve product by barcode
	if prodRec, err := l.repo.ReadByBarcode(ctx,
		l.conn.DB.WithContext(ctx.Request().Context()),
		prodRec,
		ctx.QueryParams(),
	); err != nil {
		logrus.Warnf("product.logic.SyncByImages > product.repo.ReadByBarcode error: %v", err)
		return err
	} else if err := l.SetImages(ctx, req, prodRec, imgs); err != nil {
		logrus.Warnf("product.logic.SyncByImages > product.SetImage error: %v", err)
		return err
	} else if err := l.repo.UpdateImages(ctx,
		l.conn.DB.WithContext(ctx.Request().Context()),
		prodRec.LocalProduct,
		nil); err != nil {
		logrus.Warnf("product.logic.SyncByImages > product.repo.Update error: %v", err)
		return err
	} else {
		logrus.Infof("%s Updated", prodRec.Barcode)
		return nil
	}
}

func (l *logic) FindOrCreateProductGroup(ctx *context.Context, reqLProductGroup *LocalProductGroup) (*LocalProductGroup, error) {
	// 1) First or create the product group
	if tx := l.conn.DB.WithContext(ctx.Request().Context()); tx == nil {
		return nil, nil
	} else if slprdgrp, err := l.repo.FirstOrCreateLocalProuctGroup(ctx, tx, reqLProductGroup); err != nil {
		// 2) First or create product records related to the group
		return nil, err
	} else {
		return slprdgrp, nil
	}
}

func (l *logic) Save(ctx *context.Context, reqPrdRec *ProductRecord) (prdRec *ProductRecord, err error) {
	var (
		topNodes models.Nodes
		// productsResp = models.ProductsResponse{}
		// isChanged    = false
	)

	logrus.Infof("Running task for product => %v", reqPrdRec)

	if reqPrdRec.CategorySlug != "" {
		if catRecs, err := l.catRepo.ReadCategoryRecords(ctx,
			l.conn.DB.WithContext(ctx.Request().Context()),
			url.Values{
				"slug": []string{reqPrdRec.CategorySlug},
			},
		); err != nil {
			return nil, err
		} else if len(catRecs) != 1 {
			logrus.WithFields(logrus.Fields{
				"slug":             reqPrdRec.CategorySlug,
				"product_record":   reqPrdRec,
				"category_records": catRecs,
			}).Errorln("no category record found")
			return nil, models.ErrNotFound
		} else if catRec := catRecs[0]; catRec == nil {
			// break
		} else {
			reqPrdRec.LocalCategory = catRec.LocalCategory
			topNodes = append(topNodes, reqPrdRec.LocalCategory.Category.Nodes...)
		}
	}

	lprd := reqPrdRec.LocalProduct
	rprd := lprd.Product

	for _, topNode := range topNodes {
		found := false
		for _, rprdNode := range rprd.Nodes {
			if rprdNode.ParentID != nil && *rprdNode.ParentID == topNode.ID {
				found = true
				break
			}
		}

		if !found {
			nodeSlug := fmt.Sprintf("%s-%s", topNode.Slug, rprd.Slug)
			rprd.Nodes = append(rprd.Nodes, &models.Node{
				ParentID: &topNode.ID,
				System:   new(bool),
				Alias:    nodeSlug,
				Slug:     nodeSlug,
			})
			// isChanged = true
		}
	}

	if prdRec, err = l.repo.ReadByBarcode(ctx,
		l.conn.DB.WithContext(ctx.Request().Context()),
		reqPrdRec,
		ctx.QueryParams(),
	); err != nil {
		return nil, err
	}

	// TODO: Equal remote product with local
	// isChanged = true

	// // Write product
	// if !isChanged {
	// 	return prdRec, nil
	// } else if resp, err := l.client.R().
	// 	SetBody(rprd).
	// 	SetResult(&productsResp).
	// 	SetError(&productsResp).
	// 	Post("/products"); err != nil {
	// 	logrus.Info(err)
	// 	return nil, err
	// } else if !resp.IsSuccess() {
	// 	return nil, fmt.Errorf("write product (%s) response error %s", rprd.Slug, resp.Status())
	// } else if prods := productsResp.Data.Products; len(prods) == 0 || prods[0] == nil {
	// 	return nil, ErrRemoteProductNotFound
	// } else if resultProd := prods[0]; resultProd == nil {
	// 	return nil, ErrRemoteProductNotFound
	// 	// Write uploaded files into the database
	// } else if tx := l.conn.DB.WithContext(ctx.Request().Context()); tx == nil {
	// 	return nil, nil
	// } else {
	// 	rprd = resultProd
	// 	lprd.Product = rprd
	// 	lprd.ProductID = rprd.ID
	// 	if err = l.repo.Save(ctx, tx, LocalProducts{lprd}); err != nil {
	// 		logrus.Infof("writing file %v in db failed %v", lprd, err)
	// 		return nil, err
	// 	}
	// }

	return prdRec, nil
}

func (l *logic) SyncByImages(ctx *context.Context, req *SyncByImageRequest, filters url.Values) (*simscheme.Document, error) {
	var (
		err              error
		offset           = 0
		limit            = 1000
		batchSize        = 10
		mimages          image.MapImages
		mapBarcodeImages = make(map[string]image.LocalImages)
		prodRecDoc       = simscheme.
					GetSchema().
					AddNewDocumentWithType(&ProductRecord{})
	)

	// Get barcodes from image's title if synced and there aren't any related product
	for {
		if mimages, err = l.repo.ReadImagesWithoutProduct(
			ctx,
			l.conn.DB.WithContext(ctx.Request().Context()),
			url.Values{
				"limit":  []string{cast.ToString(limit)},
				"offset": []string{cast.ToString(offset)},
			},
		); err != nil {
			return prodRecDoc, err
		} else if len(mimages) == 0 {
			break
		}

		for _, img := range mimages {
			// TODO check barcode pattern
			mapBarcodeImages[img.Image.Title] = append(mapBarcodeImages[img.Image.Title], img)
		}

		pool := pond.New(batchSize, 0)

		for barcode, imgs := range mapBarcodeImages {
			var (
				prodRec = &ProductRecord{Barcode: barcode}
			)

			pool.Submit(func() {
				if err := l.UpdateImage(ctx, req, imgs, prodRec); err != nil {
					return
				}
			})

			prodRecDoc.AddNode(prodRec)
		}

		pool.StopAndWait()
		offset += limit
	}

	return prodRecDoc, nil
}

func (l *logic) SyncBySpreadSheets(ctx *context.Context) (*simscheme.Document, error) {
	var (
		prodRecDoc = simscheme.
				GetSchema().
				AddNewDocumentWithType(&ProductRecord{})

		batchSize = 10
		pool      = pond.New(batchSize, 0)
	)

	if req, err := context.GetRequestModel[*SyncBySpreadSheetsRequest](ctx); err != nil {
		return prodRecDoc, err
	} else if req.ProductHeaderMap.Barcode == "" {
		return nil, nil
	} else if csvFiles, err := excel.LoadExcels(ctx, req.Root, req.MaxDepth); err != nil {
		return nil, err
		// Make ProductNodes from the files
	} else if err := excel.FromFiles(
		csvFiles,
		req.Offset,
		func(header map[string]int, rec []string) {
			var (
				prodRec = &ProductRecord{
					Barcode: rec[header[req.ProductHeaderMap.Barcode]],
					Title:   rec[header[req.ProductHeaderMap.Title]],
				}
			)

			if req.ProductHeaderMap.CategorySlug != "" {
				prodRec.CategorySlug = slug.Make(rec[header[req.ProductHeaderMap.CategorySlug]])
			}

			pool.Submit(func() {
				if err := l.UpdateNodes(ctx, prodRec); err != nil {
					return
				}
			})

			prodRecDoc.AddNode(prodRec)
		},
	); err != nil {
		return nil, err
	} else {
		pool.StopAndWait()
		return prodRecDoc, nil
	}
}

func (l *logic) SetImages(_ *context.Context, req *SyncByImageRequest, rec *ProductRecord, limages image.LocalImages) error {
	err := ErrProductNoChange

	lprod := rec.LocalProduct
	rlprod := lprod.Product

	if req.ReplaceGallery &&
		rlprod != nil &&
		len(rlprod.Images) > 0 {
		clear(rlprod.Images)
		rlprod.Images = nil
	}

	for _, limg := range limages {
		// TODO if image is cover then set cover is true else add to gallery
		if (req.ReplaceCover && !lprod.CoverSet) ||
			((!lprod.CoverID.IsValid() && lprod.Cover == nil) &&
				(!req.IgnoreCoverIfEmpty && !rlprod.CoverID.IsValid())) {
			// product.Cover = &ProductImage{
			// 	ImageID:   img.ID,
			// 	Image:     img,
			// 	ProductID: product.ID,
			// }
			// product.CoverID, _ = img.ID.ToNullPID()
			lprod.Cover = limg
			lprod.CoverID, _ = limg.ID.ToNullPID()
			rlprod.Cover = limg.Image
			rlprod.CoverID, _ = limg.ImageID.ToNullPID()
			lprod.CoverSet = true
		} else if req.ReplaceGallery || !req.IgnoreAddToGallery {
			lprod.Gallery = append(lprod.Gallery, &ProductImage{
				LocalImageID:   limg.ID,
				LocalImage:     limg,
				LocalProductID: lprod.ID,
			})
			rlprod.Images = append(rlprod.Images, &models.Imagable{
				ImageID: limg.ID,
				Image:   limg.Image,
			})
		} else {
			logrus.Infof("no chanes %s", rlprod.Barcodes)
			continue
		}

		err = nil
	}

	return err
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

	return nil, models.ErrNotFound
}
