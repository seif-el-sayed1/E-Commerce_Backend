package subcategory

import (
	"time"

	"github.com/google/uuid"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/category"
)

type SubCategory struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name string    `json:"name" gorm:"not null;index:subcategory_name"`

	CategoryID uuid.UUID         `json:"category_id" gorm:"not null;index"`
	Category   category.Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
