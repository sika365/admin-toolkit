package image

import (
	simutils "github.com/alifakhimi/simple-utils-go"

	"github.com/sika365/admin-tools/pkg/file"
)

type ScanRequest struct {
	file.ScanRequest
}

type ScanResponse struct {
	simutils.ResponseTemplate
	Data MapImages `json:"data,omitempty"`
}

type SyncRequest struct {
	*ScanRequest
	Replace bool `json:"replace,omitempty" query:"replace"`
}

type SyncResponse struct {
	simutils.ResponseTemplate
	Data MapImages `json:"data,omitempty"`
}
