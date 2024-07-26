package node

import (
	"fmt"

	"github.com/alifakhimi/simple-utils-go/simscheme"
)

type CategoryDoc *simscheme.Document
type CategoryRecord struct {
	Title string `json:"title,omitempty" gorm:"primaryKey"`
}

type ProductDoc *simscheme.Document
type ProductRecord struct {
	Barcode        string          `json:"barcode,omitempty" gorm:"primaryKey"`
	Title          string          `json:"title,omitempty"`
	Category       string          `json:"category,omitempty"`
	CategoryRecord *CategoryRecord `json:"category_record,omitempty"`
}

type NodeRecords []*NodeRecord

type NodeRecordType string

const (
	NodeRecordTypeCategory     NodeRecordType = "category"
	NodeRecordTypeProduct      NodeRecordType = "product"
	NodeRecordTypeProductGroup NodeRecordType = "product_variation"
)

type NodeRecordStatus string

const (
	NodeRecordStatusDraft   NodeRecordStatus = "draft"
	NodeRecordStatusPublish NodeRecordStatus = "publish"
)

type NodeRecord struct {
	ID            string           `json:"id,omitempty"`
	ParentID      string           `json:"parent_id,omitempty"`
	Type          NodeRecordType   `json:"type,omitempty"`
	Title         string           `json:"title,omitempty"`
	Content       string           `json:"content,omitempty"`
	Excerpt       string           `json:"excerpt,omitempty"`
	Slug          string           `json:"slug,omitempty"`
	Status        NodeRecordStatus `json:"status,omitempty"`
	Barcode       string           `json:"barcode,omitempty"`
	Path          string           `json:"path,omitempty"`
	PathSeparator string           `json:"path_separator,omitempty"`
	Parent        *NodeRecord      `json:"parent,omitempty"`
}

func ToNodeRecord(title, barcode, path string) *NodeRecord {
	return &NodeRecord{
		Title:   title,
		Barcode: barcode,
		Path:    path,
	}
}

func (nr *NodeRecord) String() string {
	return fmt.Sprintf(nr.Title)
}

func (nr *NodeRecord) Key() string {
	panic("uninmplemented")
}
