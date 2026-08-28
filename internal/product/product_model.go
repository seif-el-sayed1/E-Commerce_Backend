package product

import (
	"time"

	"github.com/google/uuid"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/subcategory"
)

type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;index:product_name"`
	Description string    `json:"description" gorm:"not null;"`

	SubCategoryID uuid.UUID               `json:"subcategory_id" gorm:"not null;index"`
	SubCategory   subcategory.SubCategory `json:"subcategory,omitempty" gorm:"foreignKey:SubCategoryID;references:ID"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
