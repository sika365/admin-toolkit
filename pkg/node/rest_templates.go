package node

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/utils/templates"

	"github.com/sika365/admin-tools/pkg/file"
)

type ScanRequest struct {
	models.NodeRequest
	file.ScanRequest
	Offset            int            `json:"offset,omitempty"`
	CategoryHeaderMap CategoryRecord `json:"category_header_map,omitempty"`
	ProductHeaderMap  ProductRecord  `json:"product_header_map,omitempty"`
}
type ScanResponse struct {
	templates.ResponseTemplate
	Data MapNodes `json:"data,omitempty"`
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
	Data MapNodes `json:"data,omitempty"`
}

type NodeSearchResponse struct {
	templates.ResponseTemplate
	Data *struct {
		Nodes models.Nodes `json:"nodes,omitempty"`
	} `json:"data,omitempty"`
}
