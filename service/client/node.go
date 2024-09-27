package client

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/context"
)

func (c *Client) GetNodeByAlias(ctx *context.Context, slug string) (node *models.Node, err error) {
	var nodeResp models.NodesResponse
	if resp, err := c.R().
		SetPathParams(map[string]string{
			"slug": slug,
		}).
		SetQueryParamsFromValues(url.Values{
			"limit":    []string{cast.ToString(1)},
			"includes": []string{"Parent", "Nodes", "Category"},
		}).
		// SetBody(&models.NodeRequest{
		// 	Node: models.Node{
		// 		Alias: "Uncategorized",
		// 		Slug: "uncategorized",
		// 	},
		// }).
		SetResult(&nodeResp).
		SetError(&nodeResp).
		Get("/nodes/{slug}"); err != nil {
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, fmt.Errorf(resp.Status())
	} else if nodes := nodeResp.Data.Nodes; len(nodes) == 0 {
		return nil, fmt.Errorf("slug(%s) not found", slug)
	} else {
		return nodes[0], nil
	}
}
