package category

import "gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"

type CategoryRecord struct {
	Title           string         `json:"title,omitempty" gorm:"primaryKey"`
	LocalCategoryID database.PID   `json:"local_category_id,omitempty"`
	LocalCategory   *LocalCategory `json:"local_category,omitempty"`
}

func (CategoryRecord) TableName() string {
	return "category_records"
}
