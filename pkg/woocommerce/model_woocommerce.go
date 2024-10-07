package woocommerce

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gorm.io/gorm"

	"github.com/sika365/admin-tools/pkg/category"
	"github.com/sika365/admin-tools/pkg/image"
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
	Posts              []*WpPost       `json:"posts,omitempty" gorm:"many2many:wp_term_relationships;foreignKey:term_taxonomy_id;joinForeignKey:term_taxonomy_id;references:id;joinReferences:object_id;"`
	Meta               []*WpTermmeta   `json:"meta,omitempty" gorm:"foreignKey:term_id"`
	Thumbnail          *WpPost         `json:"thumbnail,omitempty" gorm:"-:all"`
}

type WpTerm struct {
	TermID    uint   `json:"term_id,omitempty" gorm:"column:term_id;primaryKey"`
	Name      string `json:"name,omitempty" gorm:"column:name;"`
	Slug      string `json:"slug,omitempty" gorm:"column:slug;"`
	TermGroup int    `json:"term_group,omitempty" gorm:"column:term_group;"`
}

type WpTermmeta struct {
	MetaId    uint   `json:"meta_id,omitempty" gorm:"primaryKey;"`
	TermId    uint   `json:"term_id,omitempty" gorm:"column:term_id;"`
	MetaKey   string `json:"meta_key,omitempty" gorm:"column:meta_key;"`
	MetaValue string `json:"meta_value,omitempty" gorm:"column:meta_value;"`
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

func (tt *WpTermTaxonomy) AfterFind(tx *gorm.DB) (err error) {
	// Find thumbnail Id from termmeta
	var (
		thumbnailId uint
		thumbnail   WpPost
	)

	for _, m := range tt.Meta {
		switch m.MetaKey {
		case "thumbnail_id":
			thumbnailId = cast.ToUint(m.MetaValue)
		}
		if thumbnailId > 0 {
			break
		}
	}

	if thumbnailId == 0 {
		return nil
	}

	if err := tx.
		Preload("Meta").Preload("TermTaxonomies.Term").
		Preload("Posts").Preload("Posts.Meta").
		Where("id = ? AND post_type = ?", thumbnailId, "attachment").
		Find(&thumbnail).Error; err != nil {
		return err
	}

	if thumbnail.ID != 0 {
		tt.Thumbnail = &thumbnail
		tt.Posts = append(tt.Posts, tt.Thumbnail)
	}

	return nil
}

func (p *WpPost) GetBarcodes() (barcodes models.Barcodes) {
	for _, m := range p.Meta {
		if m.MetaKey == "_sku" && len(strings.TrimSpace(m.MetaValue)) >= 3 {
			barcodes = append(barcodes, &models.Barcode{Barcode: m.MetaValue})
		}
	}

	return barcodes
}

func (p *WpPost) GetCategoryRecords() category.CategoryRecords {
	catRecs := category.CategoryRecords{}

	for _, tt := range p.TermTaxonomies {
		if tt.Taxonomy != "product_cat" {
			continue
		}

		decodedSlug, err := url.QueryUnescape(tt.Term.Slug)
		if err != nil {
			fmt.Println("Error decoding string:", err)
			continue
		}
		catRecNode := catRecDoc.GetNode(
			&category.CategoryRecord{
				Slug: simutils.Slug(decodedSlug),
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

func (p *WpPost) IsProductGroup() bool {
	for _, subpost := range p.Posts {
		switch subpost.PostType {
		case models.PRODUCT_VARIATION_TYPE:
			return true
		}
	}
	return false
}

func (p *WpPost) GetAttachments() (cover *models.Image, gallery models.Imagables) {
	var (
		thumbnailId uint
	)

	for _, m := range p.Meta {
		switch m.MetaKey {
		case "_thumbnail_id":
			thumbnailId = cast.ToUint(m.MetaValue)
		}
		if thumbnailId > 0 {
			break
		}
	}

	// attachment
	for _, subpost := range p.Posts {
		switch subpost.PostType {
		case "attachment":
			img := &models.Image{
				Title: subpost.PostTitle,
				URL:   *database.URLFromString(subpost.Guid),
				Tags:  "[external_url]",
				Description: func() string {
					if subpost.PostContent != "" {
						return subpost.PostContent
					} else if subpost.PostExcerpt != "" {
						return subpost.PostExcerpt
					}
					return ""
				}(),
				Alias: subpost.PostName,
				Name:  subpost.PostTitle,
			}

			for _, m := range subpost.Meta {
				switch m.MetaKey {
				case "_wp_attached_file":
				case "_wp_attachment_metadata":
				case "_wp_attachment_image_alt":
					img.Title = m.MetaValue
				}
			}

			if thumbnailId == subpost.ID {
				cover = img
			} else {
				gallery = append(gallery, &models.Imagable{
					Image: img,
				})
			}
		}
	}

	return
}

func (p *WpTermTaxonomy) GetAttachments() (cover *models.Image, gallery models.Imagables) {
	var (
		thumbnailId uint
	)

	for _, m := range p.Meta {
		switch m.MetaKey {
		case "thumbnail_id":
			thumbnailId = cast.ToUint(m.MetaValue)
		}
		if thumbnailId > 0 {
			break
		}
	}

	// attachment
	for _, subpost := range p.Posts {
		switch subpost.PostType {
		case "attachment":
			img := &models.Image{
				Title: subpost.PostTitle,
				URL:   *database.URLFromString(subpost.Guid),
				Tags:  "[external_url]",
				Description: func() string {
					if subpost.PostContent != "" {
						return subpost.PostContent
					} else if subpost.PostExcerpt != "" {
						return subpost.PostExcerpt
					}
					return ""
				}(),
				Alias: subpost.PostName,
				Name:  subpost.PostTitle,
			}

			for _, m := range subpost.Meta {
				switch m.MetaKey {
				case "_wp_attached_file":
				case "_wp_attachment_metadata":
				case "_wp_attachment_image_alt":
					img.Title = m.MetaValue
				}
			}

			if thumbnailId == subpost.ID {
				cover = img
			} else {
				gallery = append(gallery, &models.Imagable{
					Image: img,
				})
			}
		}
	}

	return
}

func (p *WpPost) ToProduct(
	group *models.ProductGroup,
	cover *models.Image,
	gallery models.Imagables,
	topNodes models.Nodes,
) *models.Product {
	var (
		barcodes       = p.GetBarcodes()
		productGroupID database.NullPID
	)

	if group != nil {
		productGroupID, _ = group.ID.ToNullPID()
	}

	return &models.Product{
		LocalProduct: &models.LocalProduct{
			AppName:       p.PostTitle,
			SameNameInPos: new(bool),
			AllBarcodes:   barcodes.String(),
			Barcodes:      barcodes,
			Excerpt:       p.PostExcerpt,
			Description:   p.PostContent,
			// CoverID: cover.ID.ToNullPID(),
			Cover:          cover,
			Images:         gallery,
			ProductGroupID: productGroupID,
			Slug:           p.PostName,
			AdditionalModel: models.AdditionalModel{
				Active: simutils.SetToNilIfZeroValue(true),
			},
		},
	}
}

func (p *WpPost) ToProductRecord(catAlias string, prd *models.Product) *product.ProductRecord {
	var (
		// barcodes = p.GetBarcodes()
		barcodes = prd.LocalProduct.Barcodes
	)

	// Check `Meta`
	// ...
	// Iterate over sub-`Posts`
	// 	PostStatus: [inherit, ]
	// 	PostType: ["revision", "attachment", "product_variation"]
	// 		Revision:
	// 		Attachment:
	// 			Guid: url
	// 			Meta: [
	// 				MetaKey: "_wp_attached_file"
	// 				MetaKey: "_wp_attachment_metadata"
	// 				MetaKey: "_wp_attachment_image_alt"
	// 			]
	// 		Variation

	if len(barcodes) == 0 {
		logrus.WithFields(logrus.Fields{
			"post": p,
		}).Errorln(errors.New("no barcode available"))
		return nil
	} else if len(barcodes[0].Barcode) < 3 {
		logrus.WithFields(logrus.Fields{
			"post": p,
		}).Errorln(errors.New("barcode length is smaller than 3 characters"))
		return nil
	}

	prdRec := &product.ProductRecord{
		Barcode:      barcodes[0].Barcode,
		Title:        p.PostTitle,
		CategorySlug: catAlias,
		LocalProduct: product.FromProduct(prd),
	}

	return prdRec
}

func (p *WpPost) ToLocalProductGroup(topNodes models.Nodes) *product.LocalProductGroup {
	var (
		cover, gallery = p.GetAttachments()
	)

	// Create product group
	productGroup := &models.ProductGroup{
		CommonTableFields: models.CommonTableFields{
			Active:      simutils.SetToNilIfZeroValue(true),
			Description: p.PostContent,
			Meta: func() []byte {
				m, _ := json.Marshal(p.Meta)
				return m
			}(),
			Error: nil,
		},
		Name:    p.PostTitle,
		Slug:    p.PostName,
		Excerpt: p.PostExcerpt,
		Cover:   cover,
		Images:  gallery,
	}

	// product_variation
	for _, subpost := range p.Posts {
		switch subpost.PostType {
		case "product_variation":
			productGroup.Products = append(
				productGroup.Products,
				subpost.ToProduct(
					productGroup,
					cover,
					gallery,
					topNodes,
				),
			)
		}
	}

	return &product.LocalProductGroup{
		Slug:         p.PostName,
		Cover:        image.FromImage(cover),
		Gallery:      gallery,
		ProductGroup: productGroup,
	}
}

func (tt *WpTermTaxonomy) ToCategoryRecord() *category.CategoryRecord {
	var (
		parentAlias    simutils.Slug
		cover, gallery = tt.GetAttachments()
	)

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
				Cover:   image.FromImage(cover),
				Gallery: gallery,
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
					Cover:       cover,
					Description: tt.Description,
					Active:      simutils.SetToNilIfZeroValue[bool](true),
				},
			},
		}
	)

	return catRec
}
