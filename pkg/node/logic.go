package node

import (
	"net/url"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alifakhimi/simple-utils-go/simscheme"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/pkg/excel"
)

type Logic interface {
	Find(ctx *context.Context, req *SyncRequest, filters url.Values) (Nodes, error)
	CreateProduct(ctx *context.Context, subNodes Nodes, batchSize int) error
	Sync(ctx *context.Context, req *SyncRequest, filters url.Values) (Nodes, error)
}

type logic struct {
	conn   *simutils.DBConnection
	client *client.Client
	repo   Repo
}

func newLogic(repo Repo, conn *simutils.DBConnection, client *client.Client) (Logic, error) {
	l := &logic{
		conn:   conn,
		client: client,
		repo:   repo,
	}
	return l, nil
}

func (l *logic) Find(ctx *context.Context, req *SyncRequest, filters url.Values) (nodes Nodes, err error) {
	q := l.conn.DB.WithContext(ctx.Request().Context())

	if nodes, err := l.repo.Read(ctx, q, filters); err != nil {
		return nil, err
	} else {
		return nodes.GetValues(), nil
	}
}

func (l *logic) Sync(ctx *context.Context, req *SyncRequest, filters url.Values) (nodes Nodes, err error) {
	// var (
	// 	batchSize      = 5
	// 	filtersEncoded = filters.Encode()
	// 	pool           = pond.New(batchSize, 0)
	// )

	// // Read data
	// // ...
	// //

	// for barcode, imgs := range mapBarcodeImages {
	// 	pool.Submit(func() {
	// 		var (
	// 			filters, _ = url.ParseQuery(filtersEncoded)
	// 			conf       = config.Config()
	// 			code       = barcode
	// 			// tx           = l.conn.DB.WithContext(ctx.Request().Context())
	// 			productsResp      = ProductSearchResponse{}
	// 			updateProductResp = models.ProductsResponse{}
	// 		)
	// 		// https://sika365.com/admin/api/v1/nodes/root/products?order_by=newest&search=7899665999353&check_availability=false&search_products_in_nodes=true&search_in_node=false&search_in_sub_node=false&get_product_parents=false&search_in_reserved_quantity=false&search_in_limited_quantity=false&coverstatus=0&total=0&limit=20&offset=0&cover_status=-1&view=node&remote_pagination=false&remote_search=false&includes=Cover&includes=Nodes.Parent.Category&includes=Tags.Node.Category&includes=CategoryNodes&store_id=38&branch_id=47&stock_id=45
	// 		filters.Set("search", code)

	// 		// Is the image cover or for gallery?
	// 		// Retrieve product by barcode
	// 		if client, err := conf.GetRestyClient("sika365"); err != nil {
	// 			return
	// 		} else if resp, err := client.R().
	// 			SetQueryParamsFromValues(filters).
	// 			SetResult(&productsResp).
	// 			SetError(&productsResp).
	// 			Get("/nodes/root/products"); err != nil {
	// 			logrus.Info(err)
	// 			return
	// 		} else if !resp.IsSuccess() {
	// 			return
	// 		} else if prd, err := l.matchBarcode(ctx, req, code, productsResp.Data.ProductNodes); err != nil {
	// 			return
	// 		} else if prd, err := l.setImage(ctx, req, prd, imgs); err != nil {
	// 			return
	// 		} else if resp, err := client.R().
	// 			SetPathParams(map[string]string{
	// 				"id": prd.Product.ID.String(),
	// 			}).
	// 			SetBody(prd.Product).
	// 			SetResult(&updateProductResp).
	// 			SetError(&updateProductResp).
	// 			Put("/products/{id}"); err != nil {
	// 			return
	// 		} else if !resp.IsSuccess() {
	// 			return
	// 		} else if err := l.repo.Create(
	// 			ctx,
	// 			l.conn.DB.WithContext(ctx.Request().Context()),
	// 			Products{prd},
	// 		); err != nil {
	// 			return
	// 		} else {
	// 			logrus.Infof("%s Updated", prd.Product.LocalProduct.Barcodes)
	// 			nodes = append(nodes, prd)
	// 		}
	// 	})
	// }

	// pool.StopAndWait()

	return nodes, nil
}

func (l *logic) CreateProduct(ctx *context.Context, subNodes Nodes, batchSize int) error {
	// pool := pond.New(batchSize, 0)

	// for _, subNode := range subNodes {
	// 	pool.Submit(func() {
	// 		var (
	// 			conf      = config.Config()
	// 			nodesResp = models.NodesResponse{}
	// 		)

	// 		logrus.Infof("Running task for %v", subNode)
	// 		// Upload files
	// 		if client, err := conf.GetRestyClient("sika365"); err != nil {
	// 			return
	// 		} else if resp, err := client.R().
	// 			SetBody(subNode).
	// 			SetResult(&nodesResp).
	// 			SetError(&nodesResp).
	// 			Post("/nodes/products"); err != nil {
	// 			logrus.Info(err)
	// 			return
	// 		} else if !resp.IsSuccess() {
	// 			return
	// 		} else if nodes := nodesResp.Data.Nodes; len(nodes) == 0 || nodes[0] == nil {
	// 			return
	// 			// Write uploaded files into the database
	// 		} else if tx := l.conn.DB.WithContext(ctx.Request().Context()); tx == nil {
	// 			return
	// 		} else if err := l.repo.Create(ctx, tx, Nodes{&Node{
	// 			CommonTableFields: simutils.CommonTableFields{},
	// 			PolymorphicFields: simutils.PolymorphicFields{},
	// 			Name:              "",
	// 			Alias:             "",
	// 			Slug:              "",
	// 			ParentID:          &0,
	// 			Priority:          0,
	// 			Parent:            &Node{},
	// 			SubNodes:          []*Node{},
	// 			SubNodesCount:     0,
	// 		}}); err != nil {
	// 			logrus.Infof("writing node %v in db failed", subNode)
	// 			return
	// 		}
	// 	})
	// }

	// pool.StopAndWait()

	return nil
}

func (l *logic) SyncProducts(ctx *context.Context, req *SyncRequest, filters url.Values, root string, maxDepth int) (MapNodes, error) {
	prodDoc := simscheme.
		GetSchema().
		AddNewDocumentWithType(&ProductRecord{})

	if req.ScanRequest.ProductHeaderMap.Barcode == "" {
		return nil, nil
	} else if csvFiles, err := excel.LoadExcels(ctx, root, maxDepth); err != nil {
		return nil, err
		// Make ProductNodes from the files
	} else if _, mn := FromFiles(
		csvFiles,
		req.ScanRequest,
		func(header map[string]int, rec []string) {
			catRec := simscheme.
				GetSchema().
				GetDocumentByType(&CategoryRecord{}).
				GetNode(&CategoryRecord{Title: rec[header[req.ProductHeaderMap.Category]]}).
				Data.(*CategoryRecord)

			prodDoc.AddNode(&ProductRecord{
				Barcode:        rec[header[req.ProductHeaderMap.Barcode]],
				Title:          rec[header[req.ProductHeaderMap.Title]],
				Category:       rec[header[req.ProductHeaderMap.Category]],
				CategoryRecord: catRec,
			})
		},
	); len(mn) == 0 {
		return nil, nil
	} else {
		return mn, nil
	}
}
