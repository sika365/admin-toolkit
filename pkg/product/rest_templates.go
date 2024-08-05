package product

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/utils/templates"

	"github.com/alifakhimi/simple-utils-go/simscheme"
	"github.com/sika365/admin-tools/pkg/image"
)

type SpreadSheetRequest struct {
	Offset           int           `json:"offset,omitempty"`
	ProductHeaderMap ProductRecord `json:"product_header_map,omitempty"`
}

type ScanRequest struct {
	models.ProductRequest
	image.ScanRequest
	CoverNaming   string `json:"cover_naming,omitempty" query:"cover_naming"`
	GalleryNaming string `json:"gallery_naming,omitempty" query:"gallery_naming"`
	IgnoreMatch   bool   `json:"ignore_match,omitempty" query:"ignore_match"`
}
type ScanResponse struct {
	templates.ResponseTemplate
	Data MapProducts `json:"data,omitempty"`
}

type SyncByImageRequest struct {
	ScanRequest
	ReplaceCover       bool `json:"replace_cover,omitempty" query:"replace_cover"`
	ReplaceGallery     bool `json:"replace_gallery,omitempty" query:"replace_gallery"`
	IgnoreCoverIfEmpty bool `json:"ignore_cover_if_empty,omitempty" query:"ignore_cover_if_empty"`
	IgnoreAddToGallery bool `json:"ignore_add_to_gallery,omitempty" query:"ignore_add_to_gallery"`
}

type SyncBySpreadSheetsRequest struct {
	ScanRequest
	SpreadSheetRequest
	ReplaceNodes bool `json:"replace_nodes,omitempty" query:"replace_nodes"`
}
type SyncBySpreadSheetsResponse struct {
	templates.ResponseTemplate
	Data *simscheme.Document `json:"data,omitempty"`
}
type SyncByImagesResponse struct {
	templates.ResponseTemplate
	Data *simscheme.Document `json:"data,omitempty"`
}

type ProductSearchResponse struct {
	templates.ResponseTemplate
	Data *struct {
		ProductNodes models.Nodes `json:"product_nodes,omitempty"`
	} `json:"data,omitempty"`
}
