package category

import "github.com/google/uuid"

type Category struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name    string    `json:"name" gorm:"not null;unique;index:category_name"`
	CatType string    `json:"category_type" gorm:"not null;default:'product';check:type IN ('women','men')"`
}
