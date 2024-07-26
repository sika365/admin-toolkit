package category

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/utils/templates"

	"github.com/alifakhimi/simple-utils-go/simscheme"
	"github.com/sika365/admin-tools/pkg/file"
)

type SpreadSheetRequest struct {
	Offset            int            `json:"offset,omitempty"`
	CategoryHeaderMap CategoryRecord `json:"category_header_map,omitempty"`
}

type ScanRequest struct {
	models.NodeRequest
	file.ScanRequest
	SpreadSheetRequest
}
type ScanResponse struct {
	templates.ResponseTemplate
	Data Categories `json:"data,omitempty"`
}

type ConvertRequest struct {
	ScanRequest
	OutputPath   string `json:"output_path,omitempty"`
	ForceReplace bool   `json:"force_replace,omitempty"`
}
type ConvertResponse struct {
	templates.ResponseTemplate
	Data file.MapFiles `json:"data,omitempty"`
}

type SyncRequest struct {
	ScanRequest
	ReplaceNodes bool `json:"replace_nodes,omitempty" query:"replace_cover"`
}
type SyncResponse struct {
	templates.ResponseTemplate
	Data *simscheme.Document `json:"data,omitempty"`
}
