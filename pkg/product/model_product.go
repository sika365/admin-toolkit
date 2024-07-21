package product

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/pkg/image"
)

type Products []*Product

type ProductImages = []*ProductImage

type ProductImage struct {
	database.Model
	ImageID   database.PID `json:"image_id,omitempty"`
	Image     *image.Image `json:"image,omitempty"`
	ProductID database.PID `json:"owner_id,omitempty"`
}

type Gallery = ProductImages

type Product struct {
	*models.Product
	CoverID database.NullPID `json:"cover_id,omitempty"`
	Cover   *image.Image     `json:"cover,omitempty"`
	Gallery Gallery          `json:"gallery"`
}

func FromProduct(prd *models.Product) *Product {
	p := &Product{Product: prd}
	return p
}

func ToProduct(prd *Product) *models.Product {
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

func (Product) TableName() string {
	return "products"
}

func (p *Product) Key() string {
	return p.ID.String()
}

func (p *Product) ToProduct() *models.Product {
	return ToProduct(p)
}
