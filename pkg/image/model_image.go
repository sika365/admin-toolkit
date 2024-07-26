package image

import (
	"regexp"

	"github.com/sika365/admin-tools/pkg/file"
	"github.com/sika365/admin-tools/utils"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
)

type Images []*Image

func (imgs Images) Add(img *Image) Images {
	return imgs
}

type Image struct {
	models.CommonTableFields
	ImageID database.PID  `json:"image_id,omitempty"`
	Image   *models.Image `json:"image,omitempty"`
	FileID  database.PID  `json:"file_id,omitempty" gorm:"default:null"`
	File    *file.File    `json:"file,omitempty"`
}

func (Image) TableName() string {
	return "local_images"
}

func (i *Image) Hash() string {
	if i.File != nil && i.File.HashValid() {
		return i.File.Hash
	} else {
		return i.Image.Alias
	}
}

func FromFiles(files file.MapFiles, titlePattern *regexp.Regexp) (Images, MapImages) {
	m := make(MapImages)
	imgs := make(Images, 0, len(files))
	for _, f := range files {
		submatch := utils.FindStringSubmatch(titlePattern, f.Name)
		img := &Image{
			Image: &models.Image{
				Title:       submatch["title"],
				Description: submatch["description"],
			},
			FileID: f.ID,
			File:   f,
		}
		imgs = append(imgs, img)
		m.Add(img)
	}
	return imgs, m
}
