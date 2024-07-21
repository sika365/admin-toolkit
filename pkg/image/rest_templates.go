package image

import simutils "github.com/alifakhimi/simple-utils-go"

type ScanRequest struct {
	Root         string `json:"root,omitempty" query:"root"`
	MaxDepth     int    `json:"max_depth,omitempty" query:"max_depth"`
	ContentTypes string `json:"content_types,omitempty" query:"content_types"`
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
