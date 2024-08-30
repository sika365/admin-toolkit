package woocommerce

import (
	"time"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/spf13/cast"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/pkg/category"
	"github.com/sika365/admin-tools/pkg/node"
	"github.com/sika365/admin-tools/pkg/product"
)

type WpPost struct {
	ID                  uint              `json:"id,omitempty" gorm:"column:id;primaryKey;"`
	PostAuthor          uint              `json:"post_author,omitempty" gorm:"column:post_author;"`
	PostDate            time.Time         `json:"post_date,omitempty" gorm:"column:post_date;"`
	PostDateGmt         time.Time         `json:"post_date_gmt,omitempty" gorm:"column:post_date_gmt;"`
	PostContent         string            `json:"post_content,omitempty" gorm:"column:post_content;"`
	PostTitle           string            `json:"post_title,omitempty" gorm:"column:post_title;"`
	PostExcerpt         string            `json:"post_excerpt,omitempty" gorm:"column:post_excerpt;"`
	PostStatus          string            `json:"post_status,omitempty" gorm:"column:post_status;"`
	CommentStatus       string            `json:"comment_status,omitempty" gorm:"column:comment_status;"`
	PingStatus          string            `json:"ping_status,omitempty" gorm:"column:ping_status;"`
	PostPassword        string            `json:"post_password,omitempty" gorm:"column:post_password;"`
	PostName            string            `json:"post_name,omitempty" gorm:"column:post_name;"`
	ToPing              string            `json:"to_ping,omitempty" gorm:"column:to_ping;"`
	Pinged              string            `json:"pinged,omitempty" gorm:"column:pinged;"`
	PostModified        time.Time         `json:"post_modified,omitempty" gorm:"column:post_modified;"`
	PostModifiedGmt     time.Time         `json:"post_modified_gmt,omitempty" gorm:"column:post_modified_gmt;"`
	PostContentFiltered string            `json:"post_content_filtered,omitempty" gorm:"column:post_content_filtered;"`
	PostParent          uint              `json:"post_parent,omitempty" gorm:"column:post_parent;"`
	Guid                string            `json:"guid,omitempty" gorm:"column:guid;"`
	MenuOrder           int               `json:"menu_order,omitempty" gorm:"column:menu_order;"`
	PostType            string            `json:"post_type,omitempty" gorm:"column:post_type;"`
	PostMimeType        string            `json:"post_mime_type,omitempty" gorm:"column:post_mime_type;"`
	CommentCount        int               `json:"comment_count,omitempty" gorm:"column:comment_count;"`
	Meta                []*WpPostmeta     `json:"meta,omitempty" gorm:"foreignKey:post_id;"`
	Parent              *WpPost           `json:"parent,omitempty" gorm:"foreignKey:post_parent;"`
	Posts               []*WpPost         `json:"posts,omitempty" gorm:"foreignKey:post_parent;"`
	TermTaxonomies      []*WpTermTaxonomy `json:"term_taxonomies,omitempty" gorm:"many2many:wp_term_relationships;foreignKey:id;joinForeignKey:object_id;references:term_taxonomy_id;joinReferences:term_taxonomy_id;"`
}

type WpPostmeta struct {
	MetaId    uint   `json:"meta_id,omitempty" gorm:"primaryKey;"`
	PostId    uint   `json:"post_id,omitempty" gorm:"column:post_id;"`
	MetaKey   string `json:"meta_key,omitempty" gorm:"column:meta_key;"`
	MetaValue string `json:"meta_value,omitempty" gorm:"column:meta_value;"`
}

type WpTermRelationship struct {
	ObjectId       uint `json:"object_id,omitempty" gorm:"column:object_id;primaryKey;"`
	TermTaxonomyId uint `json:"term_taxonomy_id,omitempty" gorm:"column:term_taxonomy_id;primaryKey;"`
	TermOrder      int  `json:"term_order,omitempty" gorm:"column:term_order;"`
}

type WpTermTaxonomy struct {
	TermTaxonomyId     uint            `json:"term_taxonomy_id,omitempty" gorm:"column:term_taxonomy_id;primaryKey;"`
	TermId             uint            `json:"term_id,omitempty" gorm:"column:term_id;"`
	Taxonomy           string          `json:"taxonomy,omitempty" gorm:"column:taxonomy;"`
	Description        string          `json:"description,omitempty" gorm:"column:description;"`
	Parent             uint            `json:"parent,omitempty" gorm:"column:parent;"`
	Count              int             `json:"count,omitempty" gorm:"column:count;"`
	Term               *WpTerm         `json:"term,omitempty" gorm:"references:term_id;foreignKey:term_id"`
	ParentTermTaxonomy *WpTermTaxonomy `json:"parent_term_taxonomy,omitempty" gorm:"foreignKey:parent"`
}

type WpTerm struct {
	TermID    uint   `json:"term_id,omitempty" gorm:"column:term_id;primaryKey"`
	Name      string `json:"name,omitempty" gorm:"column:name;"`
	Slug      string `json:"slug,omitempty" gorm:"column:slug;"`
	TermGroup int    `json:"term_group,omitempty" gorm:"column:term_group;"`
}

func (WpPost) TableName() string {
	return "wp_posts"
}

func (WpPostmeta) TableName() string {
	return "wp_postmeta"
}

func (WpTermRelationship) TableName() string {
	return "wp_term_relationships"
}

func (WpTermTaxonomy) TableName() string {
	return "wp_term_taxonomy"
}

func (WpTerm) TableName() string {
	return "wp_terms"
}

func (p *WpPost) GetBarcode() (barcode string) {
	for _, m := range p.Meta {
		if m.MetaKey == "_sku" {
			return m.MetaValue
		}
	}

	return ""
}

func (p *WpPost) GetCategoryRecords() category.CategoryRecords {
	catRecs := category.CategoryRecords{}

	for _, tt := range p.TermTaxonomies {
		if tt.Taxonomy != "product_cat" {
			continue
		}

		catRecNode := catRecDoc.GetNode(
			&category.CategoryRecord{
				Slug: simutils.Slug(tt.Term.Slug),
			},
		)

		if catRecNode == nil {
			continue
		}

		catRec := catRecNode.Data.(*category.CategoryRecord)

		if catRec.LocalCategory != nil && len(catRec.LocalCategory.Nodes) > 0 {
			for _, node := range catRec.LocalCategory.Nodes {
				if len(node.SubNodes) == 0 {
					catRecs = append(catRecs, catRec)
				}
			}
		}
	}

	return catRecs
}

func (p *WpPost) ToProductRecord() *product.ProductRecord {
	var (
		barcode         = p.GetBarcode()
		categoryRecords = p.GetCategoryRecords()
	)

	if len(categoryRecords) == 0 || barcode == "" {
		return nil
	}

	catRec := categoryRecords[0]
	prdRec := &product.ProductRecord{
		Barcode:       barcode,
		Title:         p.PostTitle,
		CategoryAlias: string(catRec.Slug),
		LocalCategory: catRec.LocalCategory,
	}

	return prdRec
}

func (tt *WpTermTaxonomy) ToCategoryRecord() *category.CategoryRecord {
	parentAlias := simutils.Slug("")
	if tt.Parent > 0 {
		parentAlias = simutils.MakeSlug(cast.ToString(tt.Parent))
	}

	var (
		catSlug = simutils.MakeSlug(tt.Term.Slug)
		catRec  = &category.CategoryRecord{
			Title: tt.Term.Name,
			Slug:  catSlug,
			LocalCategory: &category.LocalCategory{
				Title:   tt.Term.Name,
				Alias:   catSlug,
				Slug:    catSlug,
				Content: tt.Description,
				Nodes: node.LocalNodes{
					{
						CommonTableFields: simutils.CommonTableFields{
							Description: tt.Description,
						},
						Name:        tt.Term.Name,
						Alias:       simutils.MakeSlug(cast.ToString(tt.TermTaxonomyId)),
						Slug:        catSlug,
						Priority:    0,
						ParentAlias: parentAlias,
					},
				},
				Category: &models.Category{
					Title:       tt.Term.Name,
					Slug:        catSlug.ToString(),
					Alias:       catSlug.ToString(),
					Description: tt.Description,
					Active:      simutils.SetToNilIfZeroValue[bool](true),
				},
			},
		}
	)

	return catRec
}
