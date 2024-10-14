package client

import (
	"fmt"
	"net/url"

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
		// log
		logEntry = logrus.
				WithContext(ctx.Request().Context()).
				WithFields(logrus.Fields{
				"GET":     "/products",
				"filters": cfilters,
				"result":  &productsResp,
			})
	)

	qp := "check_availability=false&search_products_in_nodes=false&search_in_node=false&search_in_sub_node=false&get_product_parents=false&search_in_reserved_quantity=false&search_in_limited_quantity=false&cover_status=0&check_products_in_nodes=false&remote_pagination=true&remote_search=true&includes=Cover&includes=Nodes.Parent.Category&includes=Tags.Node.Category&includes=Nodes&includes=CategoryNodes"
	cfilters, _ = url.ParseQuery(cfilters.Encode() + "&" + qp)
	cfilters.Set("search", barcode)

	if resp, err := c.R().
		// SetPathParams(map[string]string{
		// 	"node_id": "root",
		// }).
		SetQueryParamsFromValues(cfilters).
		SetResult(&productsResp).
		SetError(&productsResp).
		Get("/products"); err != nil {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if !resp.IsSuccess() {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln()
		return nil, fmt.Errorf(resp.Status())
	} else if products := productsResp.Data.Products; len(products) == 0 {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
			"products": products,
		}).Errorln(models.ErrNotFound)
		return nil, models.ErrNotFound
	} else if prods, err := c.MatchBarcode(barcode, products); err != nil {
		logEntry.WithFields(logrus.Fields{
			"fn":       "Product.MatchBarcode",
			"products": prods,
		}).Errorln(err)
		return nil, err
	} else {
		return prods, nil
	}
}

func (c *Client) CreateProduct(ctx *context.Context, rprd *models.Product) (*models.Product, error) {
	var (
		productsResp = models.ProductsResponse{}
		// log
		logEntry = logrus.
				WithContext(ctx.Request().Context()).
				WithFields(logrus.Fields{
				"POST":   "/products",
				"body":   rprd,
				"result": &productsResp,
			})
	)

	if resp, err := c.R().
		SetBody(rprd).
		SetResult(&productsResp).
		SetError(&productsResp).
		Post("/products"); err != nil {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if !resp.IsSuccess() {
		err = fmt.Errorf("write product (%s) response error %s", rprd.Slug, resp.Status())
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if prods := productsResp.Data.Products; len(prods) == 0 || prods[0] == nil {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
			"products": prods,
		}).Errorln(models.ErrNotFound)
		return nil, models.ErrNotFound
	} else if resultProd := prods[0]; resultProd == nil {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
			"product":  resultProd,
		}).Errorln(models.ErrNotFound)
		return nil, models.ErrNotFound
	} else {
		if !database.IsValid(resultProd.LocalProduct.StoreID) {
			resultProd.LocalProduct.StoreID = resultProd.ProductStock.StoreID
		}
		return resultProd, nil
	}
}

func (c *Client) PutProduct(ctx *context.Context, prd *models.Product) (*models.Product, error) {
	var (
		updatedProductResp = models.ProductsResponse{}
		// log
		logEntry = logrus.
				WithContext(ctx.Request().Context()).
				WithFields(logrus.Fields{
				"PUT":    fmt.Sprintf("/products/{%s}", prd.ID),
				"body":   prd,
				"result": &updatedProductResp,
			})
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
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln()
		return nil, fmt.Errorf(resp.Status())
	} else if products := updatedProductResp.Data.Products; len(products) == 0 {
		logEntry.WithFields(logrus.Fields{
			"products": products,
			"response": resp.Status(),
		}).Errorln(models.ErrNotFound)
		return nil, models.ErrNotFound
	} else {
		prod := products[0]
		prd.LocalProduct.PIDModel = prd.PIDModel
		return prod, nil
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
		// log
		logEntry = logrus.WithContext(ctx.Request().Context()).
				WithFields(logrus.Fields{
				"POST":   "/nodes/products",
				"body":   &body,
				"result": &productsResp,
			})
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
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if !resp.IsSuccess() {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
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

func (c *Client) GetProductGroupBySlug(ctx *context.Context, slug string) (*models.ProductGroup, error) {
	var (
		response = models.ProductGroupResponse{}
		filters  = url.Values{
			"limit":    []string{cast.ToString(1)},
			"includes": []string{"Cover", "Imagables", "Products"},
		}
		// log
		logEntry = logrus.WithContext(ctx.Request().Context()).
				WithFields(logrus.Fields{
				"GET":     fmt.Sprintf("/product_groups/{%s}", slug),
				"filters": filters,
				"result":  &response,
			})
	)

	if resp, err := c.R().
		SetPathParams(map[string]string{
			"slug": slug,
		}).
		SetQueryParamsFromValues(filters).
		SetResult(&response).
		SetError(&response).
		Get("/product_groups/{slug}"); err != nil {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if !resp.IsSuccess() {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if prdGrps := response.Data.ProductGroups; len(prdGrps) == 0 || prdGrps[0] == nil {
		logrus.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(models.ErrNotFound)
		return nil, models.ErrNotFound
	} else {
		return prdGrps[0], nil
	}
}

func (c *Client) CreateProductGroup(ctx *context.Context, prdgrp *models.ProductGroup) (*models.ProductGroup, error) {
	var (
		response = models.ProductGroupResponse{}
		request  = models.ProductGroupRequest{ProductGroup: *prdgrp}
		// log
		logEntry = logrus.WithContext(ctx.Request().Context()).
				WithFields(logrus.Fields{
				"POST":   "/product_groups",
				"body":   &request,
				"result": &response,
			})
	)

	if resp, err := c.R().
		SetBody(&request).
		SetResult(&response).
		SetError(&response).
		Post("/product_groups"); err != nil {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if !resp.IsSuccess() {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if prdGrps := response.Data.ProductGroups; len(prdGrps) == 0 || prdGrps[0] == nil {
		if prdGrp := response.Data.ProductGroup; prdGrp != nil {
			return prdGrp, nil
		} else {
			logEntry.WithFields(logrus.Fields{
				"response": resp.Status(),
			}).Errorln(models.ErrNotFound)
			return nil, models.ErrNotFound
		}
	} else {
		return prdGrps[0], nil
	}
}
