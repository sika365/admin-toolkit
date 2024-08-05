package product

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/pkg/image"
)

type LocalProducts []*LocalProduct

type ProductImages = []*ProductImage

type ProductImage struct {
	database.Model
	LocalImageID   database.PID      `json:"local_image_id,omitempty"`
	LocalImage     *image.LocalImage `json:"local_image,omitempty"`
	LocalProductID database.PID      `json:"local_product_id,omitempty"`
}

type Gallery = ProductImages

type LocalProduct struct {
	models.CommonTableFields
	CoverID   database.NullPID  `json:"cover_id,omitempty"`
	ProductID database.PID      `json:"product_id,omitempty"`
	Cover     *image.LocalImage `json:"cover,omitempty"`
	Gallery   Gallery           `json:"gallery"`
	Product   *models.Product   `json:"product,omitempty"`
}

func FromProduct(prd *models.Product) *LocalProduct {
	p := &LocalProduct{Product: prd}
	return p
}

func ToProduct(prd *LocalProduct) *models.Product {
	product := prd.Product
	product.CoverID = prd.CoverID
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

	return product
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
	if rprod := p.Product; rprod == nil {
		return ErrRemoteLocalProductNotFound
	} else if rlprod := rprod.LocalProduct; rlprod == nil {
		return ErrRemoteLocalProductNotFound
	} else {
		rlprod.Nodes = models.Nodes{}
		return nil
	}
}

func (p *LocalProduct) AddTopNodes(topNodes models.Nodes, replace bool) error {
	if replace {
		if err := p.RemoveNodes(); err != nil {
			return err
		}
	}

	if rprod := p.Product; rprod == nil {
		return ErrRemoteLocalProductNotFound
	} else if rlprod := rprod.LocalProduct; rlprod == nil {
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
			}
		}
	}
	return nil
}
