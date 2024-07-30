package client

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/utils/templates"

	"github.com/sika365/admin-tools/context"
)

type ProductSearchResponse struct {
	templates.ResponseTemplate
	Data *struct {
		ProductNodes models.Nodes `json:"product_nodes,omitempty"`
	} `json:"data,omitempty"`
}

func (c *Client) GetProductbyBarcode(ctx *context.Context, barcode string, filters url.Values) (product *models.Product, err error) {
	var (
		productsResp ProductSearchResponse
		// clone filters
		cfilters, _ = url.ParseQuery(filters.Encode())
	)

	cfilters.Set("search", barcode)

	if resp, err := c.R().
		SetPathParams(map[string]string{
			"node_id": "root",
		}).
		SetQueryParamsFromValues(cfilters).
		SetResult(&productsResp).
		SetError(&productsResp).
		Get("/nodes/{node_id}/products"); err != nil {
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, fmt.Errorf(resp.Status())
	} else if products := productsResp.Data.ProductNodes; len(products) == 0 {
		return nil, models.ErrNotFound
	} else if prod, err := c.MatchBarcode(barcode, products); err != nil {
		return nil, err
	} else {
		return prod, nil
	}
}

func (c *Client) PutProduct(ctx *context.Context, prd *models.Product) (*models.Product, error) {
	var (
		updatedProductResp = models.ProductsResponse{}
	)

	if resp, err := c.R().
		SetPathParams(map[string]string{
			"id": prd.ID.String(),
		}).
		SetBody(prd).
		SetResult(&updatedProductResp).
		SetError(&updatedProductResp).
		Put("/products/{id}"); err != nil {
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, fmt.Errorf(resp.Status())
	} else if products := updatedProductResp.Data.Products; len(products) == 0 {
		return nil, models.ErrNotFound
	} else {
		return products[0], nil
	}
}

func (c *Client) AddToNodes(ctx *context.Context, prd *models.Product, nodeIDs database.PIDs) (prodNodes models.Nodes, err error) {
	if len(nodeIDs) == 0 {
		return models.Nodes{}, nil
	}
	var (
		productsResp ProductSearchResponse
		body         = struct {
			NodeIDs   database.PIDs             `json:"node_ids,omitempty"`
			Nodes     map[database.PID][]string `json:"nodes,omitempty"`
			RemoteIDs []string                  `json:"remote_ids,omitempty"`
		}{
			NodeIDs:   nodeIDs,
			Nodes:     map[database.PID][]string{},
			RemoteIDs: []string{},
		}
	)

	for _, nid := range nodeIDs {
		rprodID := prd.ProductStock.RemoteID
		body.Nodes[nid] = append(body.Nodes[nid], rprodID)
		body.RemoteIDs = append(body.RemoteIDs, rprodID)
	}

	if resp, err := c.R().
		SetBody(body).
		SetResult(&productsResp).
		SetError(&productsResp).
		Post("/nodes/products"); err != nil {
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, fmt.Errorf(resp.Status())
		/* } else if prodNodes = productsResp.Data.ProductNodes; len(prodNodes) == 0 {
		return nil, models.ErrNotFound */
	} else {
		return prodNodes, nil
	}
}

func (c *Client) MatchBarcode(barcode string, productNodes models.Nodes) (*models.Product, error) {
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
