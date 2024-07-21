package product

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/utils/templates"

	"github.com/sika365/admin-tools/pkg/image"
)

type ScanRequest struct {
	models.ProductRequest
	*image.ScanRequest
	CoverNaming   string `json:"cover_naming,omitempty" query:"cover_naming"`
	GalleryNaming string `json:"gallery_naming,omitempty" query:"gallery_naming"`
	IgnoreMatch   bool   `json:"ignore_match,omitempty" query:"ignore_match"`
}
type ScanResponse struct {
	templates.ResponseTemplate
	Data MapProducts `json:"data,omitempty"`
}

type SyncRequest struct {
	*ScanRequest
	ReplaceCover       bool `json:"replace_cover,omitempty" query:"replace_cover"`
	ReplaceGallery     bool `json:"replace_gallery,omitempty" query:"replace_gallery"`
	IgnoreCoverIfEmpty bool `json:"ignore_cover_if_empty,omitempty" query:"ignore_cover_if_empty"`
	IgnoreAddToGallery bool `json:"ignore_add_to_gallery,omitempty" query:"ignore_add_to_gallery"`
}
type SyncResponse struct {
	templates.ResponseTemplate
	Data MapProducts `json:"data,omitempty"`
}

type ProductSearchResponse struct {
	templates.ResponseTemplate
	Data *struct {
		ProductNodes models.Nodes `json:"product_nodes,omitempty"`
	} `json:"data,omitempty"`
}
