package category

import (
	simutils "github.com/alifakhimi/simple-utils-go"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/pkg/image"
	"github.com/sika365/admin-tools/pkg/node"
)

type Categories []*Category

type CategoryNode struct{}

type Category struct {
	simutils.CommonTableFields
	Title    string           `json:"title,omitempty"`
	Alias    string           `json:"alias,omitempty" gorm:"index"`
	Slug     string           `json:"slug,omitempty" gorm:"index"`
	Content  string           `json:"content,omitempty"`
	CoverID  database.NullPID `json:"cover_id,omitempty"`
	Cover    *image.Image     `json:"cover,omitempty"`
	Nodes    node.Nodes       `json:"nodes,omitempty" gorm:"polymorphic:Owner"`
	Category *models.Category `json:"category,omitempty"`
}

func (Category) TableName() string {
	return "local_categories"
}

func (n *Category) Key() string {
	return n.ID.String()
}
