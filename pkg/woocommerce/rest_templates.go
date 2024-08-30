package woocommerce

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/utils/templates"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alifakhimi/simple-utils-go/simscheme"
	"github.com/sika365/admin-tools/pkg/image"
)

type ScanRequest struct {
	models.ProductRequest
	image.ScanRequest
	CoverNaming   string `json:"cover_naming,omitempty" query:"cover_naming"`
	GalleryNaming string `json:"gallery_naming,omitempty" query:"gallery_naming"`
	IgnoreMatch   bool   `json:"ignore_match,omitempty" query:"ignore_match"`
}

type SyncRequest struct {
	ScanRequest
	DatabaseConfig simutils.DBConfig `json:"database_config,omitempty"`
	ReplaceNodes   bool              `json:"replace_nodes,omitempty" query:"replace_nodes"`
}
type SyncResponse struct {
	templates.ResponseTemplate
	Data *simscheme.Document `json:"data,omitempty"`
}
