package client

import (
	"errors"

	"github.com/alifakhimi/simple-utils-go/simrest"
	"github.com/go-resty/resty/v2"

	"github.com/sika365/admin-tools/config"
)

var (
	ErrClientNoFound = errors.New("no client info found")
)

type Client struct {
	*simrest.Client
}

var c *Client

func New(config *config.ServiceConfig) (*Client, error) {
	if c != nil {
		return c, nil
	} else if sika365Client, err := config.GetClient("sika365"); err != nil {
		return nil, err
	} else {
		c = &Client{
			Client: sika365Client,
		}
		return c, nil
	}
}

func R() *resty.Request { return c.R() }
func (c *Client) R() *resty.Request {
	return c.Client.Client.R()
}
