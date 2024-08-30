package category

import (
	simutils "github.com/alifakhimi/simple-utils-go"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
)

type CategoryRecords []*CategoryRecord

type CategoryRecord struct {
	Slug            simutils.Slug  `json:"slug,omitempty" gorm:"primaryKey"`
	Title           string         `json:"title,omitempty"`
	LocalCategoryID database.PID   `json:"local_category_id,omitempty"`
	LocalCategory   *LocalCategory `json:"local_category,omitempty"`
}

func (CategoryRecord) TableName() string {
	return "category_records"
}

func (crs CategoryRecords) ToNested() (nested CategoryRecords) {
	mns := make(map[simutils.Slug]*CategoryRecord)

	for _, cr := range crs {
		if cr.LocalCategory != nil && len(cr.LocalCategory.Nodes) > 0 {
			mns[cr.LocalCategory.Nodes[0].Alias] = cr
			for _, node := range cr.LocalCategory.Nodes {
				if node.Parent == nil && len(node.ParentAlias) == 0 || node.ParentAlias == "0" {
					nested = append(nested, cr)
				}
			}
		}
	}

	for _, cr := range crs {
		if cr.LocalCategory != nil {
			for _, node := range cr.LocalCategory.Nodes {
				if node.ParentAlias.IsValid() {
					if crParent, ok := mns[node.ParentAlias]; ok &&
						crParent.LocalCategory != nil &&
						len(crParent.LocalCategory.Nodes) > 0 {
						parentNode := crParent.LocalCategory.Nodes[0]
						// node.Parent = crParent.LocalCategory.Nodes[0]
						parentNode.SubNodes = append(parentNode.SubNodes, node)
						parentNode.SubNodesCount++
					}
				}
			}
		}
	}

	return nested
}
