package admin

import (
	"strconv"
	"strings"
	"time"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"

	"github.com/golang-jwt/jwt/v5"
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
	IsActive     bool `json:"is_active" gorm:"default:true"`

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

func (a *Admin) GenerateToken(db *gorm.DB) (string, time.Time, error) {
	expStr := config.Env.JWTExpiration
	days, err := strconv.Atoi(strings.TrimSuffix(expStr, "d"))
	if err != nil {
		return "", time.Time{}, err
	}
	tokenExpDate := time.Now().AddDate(0, 0, days)

	claims := jwt.MapClaims{
		"userId": a.ID.String(),
		"role":   a.Role,
		"exp":    tokenExpDate.Unix(),
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenObj.SignedString([]byte(config.Env.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	a.Token = &token
	a.TokenExpDate = &tokenExpDate
	if err := db.Save(a).Error; err != nil {
		return "", time.Time{}, err
	}

	return token, tokenExpDate, nil
}
