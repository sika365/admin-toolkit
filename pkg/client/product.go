package client

import (
	"fmt"
	"net/url"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
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

func (c *Client) GetProductsByBarcode(ctx *context.Context, barcode string, filters url.Values) (products models.Products, err error) {
	var (
		// productsResp ProductSearchResponse
		productsResp models.ProductsResponse
		// clone filters
		cfilters, _ = url.ParseQuery(filters.Encode())
	)

	cfilters.Set("search", barcode)

	if resp, err := c.R().
		// SetPathParams(map[string]string{
		// 	"node_id": "root",
		// }).
		SetQueryParamsFromValues(cfilters).
		SetResult(&productsResp).
		SetError(&productsResp).
		Get("/products"); err != nil {
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, fmt.Errorf(resp.Status())
	} else if products := productsResp.Data.Products; len(products) == 0 {
		return nil, models.ErrNotFound
	} else if prods, err := c.MatchBarcode(barcode, products); err != nil {
		return nil, err
	} else {
		return prods, nil
	}
}

func (c *Client) CreateProduct(ctx *context.Context, rprd *models.Product) (*models.Product, error) {
	var (
		productsResp = models.ProductsResponse{}
	)

	if resp, err := c.R().
		SetBody(rprd).
		SetResult(&productsResp).
		SetError(&productsResp).
		Post("/products"); err != nil {
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, fmt.Errorf("write product (%s) response error %s", rprd.Slug, resp.Status())
	} else if prods := productsResp.Data.Products; len(prods) == 0 || prods[0] == nil {
		return nil, models.ErrNotFound
	} else if resultProd := prods[0]; resultProd == nil {
		return nil, models.ErrNotFound
	} else {
		return resultProd, nil
	}
}

func (c *Client) PutProduct(ctx *context.Context, prd *models.Product) (*models.Product, error) {
	var (
		updatedProductResp = models.ProductsResponse{}
	)

	if !database.IsValid(prd.ID) {
		return prd, nil
	} else if resp, err := c.R().
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

func (c *Client) MatchBarcode(barcode string, prods models.Products) (products models.Products, err error) {
	for _, p := range prods {
		if p.LocalProduct == nil {
			continue
		}
		for _, b := range p.LocalProduct.Barcodes {
			if b.Barcode == barcode {
				products = append(products, p)
				break
			}
		}
	}

	return products, nil
}

func (c *Client) MatchBarcodeByNodes(barcode string, productNodes models.Nodes) (*models.Product, error) {
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

func (c *Client) GetProductGroupBySlug(ctx *context.Context, slug simutils.Slug) (*models.ProductGroup, error) {
	var (
		response = models.ProductGroupResponse{}
	)

	if resp, err := c.R().
		SetPathParams(map[string]string{
			"slug": string(slug),
		}).
		SetQueryParamsFromValues(url.Values{
			"limit":    []string{cast.ToString(1)},
			"includes": []string{"Cover", "Imagables", "Products"},
		}).
		SetResult(&response).
		SetError(&response).
		Get("/product_groups/{slug}"); err != nil {
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, err
	} else if prdGrps := response.Data.ProductGroups; len(prdGrps) == 0 || prdGrps[0] == nil {
		return nil, models.ErrNotFound
	} else {
		return prdGrps[0], nil
	}
}

func (c *Client) CreateProductGroup(ctx *context.Context, prdgrp *models.ProductGroup) (*models.ProductGroup, error) {
	var (
		response = models.ProductGroupResponse{}
		request  = models.ProductGroupRequest{ProductGroup: *prdgrp}
	)

	if resp, err := c.R().
		SetBody(request).
		SetResult(&response).
		SetError(&response).
		Post("/product_groups"); err != nil {
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, err
	} else if prdGrps := response.Data.ProductGroups; len(prdGrps) == 0 || prdGrps[0] == nil {
		if prdGrp := response.Data.ProductGroup; prdGrp != nil {
			return prdGrp, nil
		} else {
			return nil, models.ErrNotFound
		}
	} else {
		return prdGrps[0], nil
	}
}
