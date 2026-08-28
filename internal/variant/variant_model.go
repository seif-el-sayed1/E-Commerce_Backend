package variant

import (
	"time"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/product"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Variant struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	Color string `json:"color"`
	Size  string `json:"size"`

	Price  float64        `json:"price" gorm:"not null"`
	Stock  int            `json:"stock" gorm:"not null;default:0"`
	Images pq.StringArray `json:"images" gorm:"type:text[]"`

	HasSale        bool    `json:"has_sale" gorm:"default:false"`
	SalePercent    float64 `json:"sale_percent" gorm:"default:0"`
	PriceAfterSale float64 `json:"price_after_sale" gorm:"default:0"`

	IsAvailable bool `json:"is_available" gorm:"default:false"`

	ProductID uuid.UUID       `json:"product_id" gorm:"not null;index"`
	Product   product.Product `json:"-" gorm:"foreignKey:ProductID;references:ID"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (v *Variant) BeforeSave(tx *gorm.DB) error {
	if !v.HasSale || v.SalePercent <= 0 {
		v.PriceAfterSale = v.Price
	} else {
		v.PriceAfterSale = v.Price - (v.Price * v.SalePercent / 100)
	}

	v.IsAvailable = v.Stock > 0

	return nil
}
