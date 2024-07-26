package image

import (
	"regexp"

	"github.com/sika365/admin-tools/pkg/file"
	"github.com/sika365/admin-tools/utils"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
)

type LocalImages []*LocalImage

func (imgs LocalImages) Add(img *LocalImage) LocalImages {
	return imgs
}

type LocalImage struct {
	models.CommonTableFields
	ImageID database.PID  `json:"image_id,omitempty"`
	Image   *models.Image `json:"image,omitempty"`
	FileID  database.PID  `json:"file_id,omitempty" gorm:"default:null"`
	File    *file.File    `json:"file,omitempty"`
}

func (LocalImage) TableName() string {
	return "local_images"
}

func (i *LocalImage) Hash() string {
	if i.File != nil && i.File.HashValid() {
		return i.File.Hash
	} else {
		return i.Image.Alias
	}
}

func FromFiles(files file.MapFiles, titlePattern *regexp.Regexp) (LocalImages, MapImages) {
	m := make(MapImages)
	imgs := make(LocalImages, 0, len(files))
	for _, f := range files {
		submatch := utils.FindStringSubmatch(titlePattern, f.Name)
		img := &LocalImage{
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
