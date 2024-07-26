package client

import (
	"fmt"

	"github.com/sika365/admin-tools/context"
	"github.com/sirupsen/logrus"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
)

func (c *Client) StoreImage(ctx *context.Context, filepath string, img *models.Image) (*models.Image, error) {
	var imageResp models.ImagesResponse
	if resp, err := c.R().
		SetFile("files", filepath).
		SetFormData(map[string]string{
			"title":       img.Title,
			"description": img.Description,
		}).
		SetResult(&imageResp).
		SetError(&imageResp).
		Post("/images"); err != nil {
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, fmt.Errorf("create image failed: %v", resp.Status())
	} else if imgs := imageResp.Data.Images; len(imgs) == 0 {
		return nil, models.ErrNotFound
	} else {
		return imgs[0], nil
	}
}
