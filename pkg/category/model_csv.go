package category

import "gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"

type CategoryRecord struct {
	Title      string       `json:"title,omitempty" gorm:"primaryKey"`
	CategoryID database.PID `json:"category_id,omitempty"`
	Category   *Category    `json:"category,omitempty"`
}
