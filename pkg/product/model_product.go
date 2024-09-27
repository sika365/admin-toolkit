package product

import (
	"time"

	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/helpers"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
	"gorm.io/gorm"
)

type ViwProductStock struct {
	models.ProductStock
}

// ViwProduct represents the viw_products table in the database.
type ViwProduct struct {
	ID                  database.PID     `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt           time.Time        `json:"created_at" gorm:"column:created_at"`
	UpdatedAt           time.Time        `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt           gorm.DeletedAt   `json:"deleted_at" gorm:"index"`
	OwnerType           string           `json:"owner_type" gorm:"column:owner_type"`
	OwnerID             database.PID     `json:"owner_id" gorm:"column:owner_id"`
	StoreID             database.PID     `json:"store_id" gorm:"column:store_id"`
	BranchID            database.PID     `json:"branch_id" gorm:"column:branch_id"`
	StockID             database.PID     `json:"stock_id" gorm:"column:stock_id"`
	ProductID           database.PID     `json:"product_id" gorm:"column:product_id"`
	RemoteID            string           `json:"remote_id" gorm:"column:remote_id;default:null"`
	NameInPOS           string           `json:"name_in_pos" gorm:"column:name_in_pos;default:null"`
	ListPrice           float64          `json:"list_price" gorm:"column:list_price;default:null"`
	SalesPrice          float64          `json:"sales_price" gorm:"column:sales_price;default:null"`
	Discount            float64          `json:"discount" gorm:"column:discount;default:null"`
	Quantity            float64          `json:"quantity" gorm:"column:quantity;default:null"`
	ReservedQuantity    float64          `json:"reserved_quantity" gorm:"column:reserved_quantity;default:null"`
	LimitBuyQuantity    float64          `json:"limit_buy_quantity" gorm:"column:limit_buy_quantity;default:null"`
	Version             int              `json:"version" gorm:"column:version;default:1"`
	CreatedBy           database.PID     `json:"created_by" gorm:"column:created_by"`
	UpdatedBy           database.PID     `json:"updated_by" gorm:"column:updated_by"`
	Comment             string           `json:"comment" gorm:"column:comment"`
	Active              bool             `json:"active" gorm:"column:active;default:false"`
	Status              int              `json:"status" gorm:"column:status"`
	Meta                string           `json:"meta" gorm:"column:meta;type:json"`
	Slug                string           `json:"slug" gorm:"column:slug"`
	AppName             string           `json:"app_name" gorm:"column:app_name;default:null"`
	SameNameInPOS       bool             `json:"same_name_in_pos" gorm:"column:same_name_in_pos"`
	AllBarcodes         string           `json:"all_barcodes" gorm:"column:all_barcodes;default:null"`
	Excerpt             string           `json:"excerpt" gorm:"column:excerpt;default:null"`
	Description         string           `json:"description" gorm:"column:description;default:null"`
	CoverID             database.NullPID `json:"cover_id" gorm:"column:cover_id"`
	UnitTypeID          database.PID     `json:"unit_type_id" gorm:"column:unit_type_id"`
	ProductGroupID      database.PID     `json:"product_group_id" gorm:"column:product_group_id"`
	SpecificationNodeID database.PID     `json:"specification_node_id" gorm:"column:specification_node_id"`
	QuestionnaireNodeID database.PID     `json:"questionnaire_node_id" gorm:"column:questionnaire_node_id"`
	Type                string           `json:"type" gorm:"column:type"`
	SearchText          string           `json:"search_text" gorm:"column:search_text"`
	Barcode             string           `json:"barcode" gorm:"column:barcode"`
	Name                string           `json:"name" gorm:"column:name"`
	Available           bool             `json:"available" gorm:"column:available"`
	RemoteUnitTypeID    string           `json:"remote_unit_type_id" gorm:"column:remote_unit_type_id"`
	RemoteUnitType      string           `json:"remote_unit_type" gorm:"column:remote_unit_type"`
	BrandID             string           `json:"brand_id" gorm:"column:brand_id"`
	BrandName           string           `json:"brand_name" gorm:"column:brand_name"`
	RemoteStockID       string           `json:"remote_stock_id" gorm:"column:remote_stock_id"`
	RemoteStock         string           `json:"remote_stock" gorm:"column:remote_stock"`
	CurrentDiscount     float64          `json:"current_discount" gorm:"column:current_discount"`
	Taxable             bool             `json:"taxable" gorm:"column:taxable"`
	VATPercent          float64          `json:"vat_percent" gorm:"column:vat_percent"`
	TollPercent         float64          `json:"toll_percent" gorm:"column:toll_percent"`
}

// TableName overrides the default table name used by GORM.
func (ViwProduct) TableName() string {
	return "viw_products"
}

func (vp *ViwProduct) ToProduct() *models.Product {
	var p models.Product
	if err := helpers.JSONCopy(vp, &p); err != nil {
		return nil
	} else if err := helpers.JSONCopy(vp, &p.ProductStock); err != nil {
		return nil
	}

	return &p
}

func ToViwProduct(p *models.Product) *ViwProduct {
	var vp ViwProduct
	if err := helpers.JSONCopy(p, &vp); err != nil {
		return nil
	}
	return &vp
}
