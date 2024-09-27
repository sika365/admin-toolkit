package product

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/pkg/image"
)

type LocalProducts []*LocalProduct
type LocalProductGroups []*LocalProductGroup
type ProductImages = []*ProductImage

type ProductImage struct {
	database.Model
	LocalImageID   database.PID      `json:"local_image_id,omitempty"`
	LocalImage     *image.LocalImage `json:"local_image,omitempty"`
	LocalProductID database.PID      `json:"local_product_id,omitempty"`
}

type Gallery = ProductImages

type LocalProductGroup struct {
	models.CommonTableFields
	Slug           string               `json:"slug,omitempty" query:"slug" param:"slug" sim:"primaryKey;"`
	CoverID        database.NullPID     `json:"cover_id,omitempty"`
	Cover          *image.LocalImage    `json:"cover,omitempty"`
	Gallery        models.Imagables     `json:"gallery" gorm:"polymorphic:Owner;"`
	ProductGroupID database.PID         `json:"product_group_id,omitempty"`
	ProductGroup   *models.ProductGroup `json:"product_group,omitempty"`
}

type LocalProduct struct {
	models.CommonTableFields
	ProductID database.PID         `json:"product_id,omitempty"`
	CoverID   database.NullPID     `json:"cover_id,omitempty"`
	Cover     *image.LocalImage    `json:"cover,omitempty"`
	Gallery   Gallery              `json:"gallery"`
	CoverSet  bool                 `json:"-" gorm:"-:all"`
	Product   *models.LocalProduct `json:"product,omitempty"`
}

func FromProduct(prd *models.Product) *LocalProduct {
	p := &LocalProduct{Product: prd.LocalProduct}
	p.Product.ID = prd.ID
	p.Product.StoreID = prd.StoreID
	return p
}

func ToProduct(prd *LocalProduct) *models.Product {
	product := prd.Product
	product.CoverID = database.ToNullPID(prd.Cover.ImageID)
	product.Images = make(models.Imagables, 0, len(prd.Gallery))

	for _, img := range prd.Gallery {
		found := false
		for _, i := range product.Images {
			if i.ImageID == img.ID {
				found = true
				break
			}
		}
		if found {
			continue
		}

		product.Images = append(product.Images, &models.Imagable{ImageID: img.ID})
	}

	p := &models.Product{
		PIDModel:     models.PIDModel{ID: product.ID},
		LocalProduct: product,
	}

	if len(product.ProductStocks) > 0 {
		p.ProductStock = *product.ProductStocks[0]
	}

	return p
}

func (LocalProduct) TableName() string {
	return "local_products"
}

func (p *LocalProduct) Key() string {
	return p.ID.String()
}

func (p *LocalProduct) ToProduct() *models.Product {
	return ToProduct(p)
}

func (p *LocalProduct) RemoveNodes() error {
	if rlprod := p.Product; rlprod == nil {
		return ErrRemoteLocalProductNotFound
	} else {
		rlprod.Nodes = models.Nodes{}
		return nil
	}
}

func (p *LocalProduct) AddTopNodes(topNodes models.Nodes, replace bool) error {
	err := ErrProductNoChange

	if replace {
		if err := p.RemoveNodes(); err != nil {
			return err
		}
	}

	if rlprod := p.Product; rlprod == nil {
		return ErrRemoteLocalProductNotFound
	} else if len(rlprod.Nodes) == 0 {
		for _, tnode := range topNodes {
			rlprod.Nodes = append(rlprod.Nodes, &models.Node{
				StoreID:   tnode.StoreID,
				ParentID:  &tnode.ID,
				OwnerID:   rlprod.ID,
				OwnerType: "product",
			})
		}
	} else {
		for _, tnode := range topNodes {
			exists := false
			for _, rnode := range rlprod.Nodes {
				if rnode.ParentID != nil && database.IsValid(rnode.ParentID) && tnode.ID == *rnode.ParentID {
					exists = true
					break
				}
			}
			if !exists {
				rlprod.Nodes = append(rlprod.Nodes, &models.Node{
					StoreID:   tnode.StoreID,
					ParentID:  &tnode.ID,
					OwnerID:   rlprod.ID,
					OwnerType: "product",
				})
				err = nil
			}
		}
	}
	return err
}
