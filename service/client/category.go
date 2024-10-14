package client

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/context"
)

func (c *Client) GetCategoryByAlias(ctx *context.Context, slug string) (category *models.Category, err error) {
	var (
		categoryResp models.CategoriesResponse
		// filters
		filters = url.Values{
			"limit":    []string{cast.ToString(1)},
			"includes": []string{"Nodes.Parent"},
			"excludes": []string{"product_nodes", "current_node"},
		}
		// log
		logEntry = logrus.
				WithContext(ctx.Request().Context()).
				WithFields(logrus.Fields{
				"GET":     fmt.Sprintf("/categories/{%s}", slug),
				"filters": filters,
				"result":  &categoryResp,
			})
	)

	if resp, err := c.R().
		SetPathParams(map[string]string{
			"alias": slug,
		}).
		SetQueryParamsFromValues(filters).
		SetResult(&categoryResp).
		SetError(&categoryResp).
		Get("/categories/{alias}"); err != nil {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if !resp.IsSuccess() {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln()
		return nil, fmt.Errorf(resp.Status())
	} else if categories := categoryResp.Data.Categories; len(categories) == 0 {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(models.ErrNotFound)
		return nil, models.ErrNotFound
	} else {
		return categories[0], nil
	}
}

func (c *Client) StoreCategory(ctx *context.Context, category *models.Category, in ...*models.Node) (*models.Category, error) {
	var (
		categoryResp models.CategoriesResponse
		nodeIDs      = func() (nodeIDs database.PIDs) {
			for _, n := range in {
				nodeIDs = append(nodeIDs, n.ID)
			}
			return
		}()
		// body
		body = &models.CategoryRequest{
			AddedNodes: nodeIDs,
			Category:   *category,
		}
		// log
		logEntry = logrus.
				WithContext(ctx.Request().Context()).
				WithFields(logrus.Fields{
				"GET":    "/categories",
				"body":   body,
				"result": &categoryResp,
			})
	)

	if resp, err := c.R().
		SetBody(body).
		SetResult(&categoryResp).
		SetError(&categoryResp).
		Post("/categories"); err != nil {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, err
	} else if !resp.IsSuccess() {
		logEntry.WithFields(logrus.Fields{
			"response": resp.Status(),
		}).Errorln(err)
		return nil, fmt.Errorf("create categoriy failed: %v", resp.Status())
	} else if categories := categoryResp.Data.Categories; len(categories) == 0 {
		logEntry.WithFields(logrus.Fields{
			"categories": categories,
			"response":   resp.Status(),
		}).Errorln(models.ErrNotFound)
		return nil, models.ErrNotFound
	} else {
		return categories[0], nil
	}
}
