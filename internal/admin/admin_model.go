package admin

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FirstName string    `json:"first_name" gorm:"not null;index:admin_name"`
	LastName  string    `json:"last_name" gorm:"not null;index:admin_name"`
	Email     string    `json:"email" gorm:"not null;unique"`
	Phone     *string   `json:"phone" gorm:"unique"`
	Role      string    `json:"role" gorm:"not null;default:'admin';check:role IN ('super_admin','admin')"`

	IsSuperAdmin bool `json:"is_super_admin" gorm:"default:false"`
	IsVerified   bool `json:"is_verified" gorm:"default:false"`
	IsBlocked    bool `json:"is_blocked" gorm:"default:false"`
	IsDeleted    bool `json:"is_deleted" gorm:"default:false"`

	Password               string     `json:"password" gorm:"not null"`
	PasswordChangedAt      *time.Time `json:"password_changed_at"`
	PasswordResetToken     *string    `json:"password_reset_token"`
	PasswordResetExpiresAt *time.Time `json:"password_reset_expires_at"`

	VerificationToken    *string    `json:"verification_token"`
	VerificationTokenExp *time.Time `json:"verification_token_expires_at"`

	Token        *string    `json:"token"`
	TokenExpDate *time.Time `json:"token_expires_at"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
