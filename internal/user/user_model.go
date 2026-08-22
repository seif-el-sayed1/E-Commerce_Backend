package user

import (
	"strconv"
	"strings"
	"time"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FirstName string    `json:"first_name" gorm:"not null;index:user_name"`
	LastName  string    `json:"last_name" gorm:"not null;index:user_name"`
	Email     string    `json:"email" gorm:"not null;unique"`
	Phone     *string   `json:"phone" gorm:"unique"`
	Role      string    `json:"role" gorm:"not null;default:'user'"`

	IsVerified bool `json:"is_verified" gorm:"default:false"`
	IsBlocked  bool `json:"is_blocked" gorm:"default:false"`
	IsDeleted  bool `json:"is_deleted" gorm:"default:false"`
	IsActive   bool `json:"is_active" gorm:"default:true"`

	Password               string     `json:"password" gorm:"not null"`
	PasswordChangedAt      *time.Time `json:"password_changed_at"`
	PasswordResetToken     *string    `json:"password_reset_token"`
	PasswordResetExpiresAt *time.Time `json:"password_reset_expires_at"`

	VerificationCode         *string    `json:"verification_code"`
	VerificationCodeExp      *time.Time `json:"verification_code_expires_at"`
	VerificationCodeVerified bool       `json:"verification_code_verified"`
	VerificationToken        *string    `json:"verification_token"`
	VerificationTokenExp     *time.Time `json:"verification_token_expires_at"`

	Token        *string    `json:"token"`
	TokenExpDate *time.Time `json:"token_expires_at"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (a User) GetToken() *string {
	return a.Token
}

func (a User) GetIsActive() bool {
	return a.IsActive
}

func (a User) GetIsBlocked() bool {
	return a.IsBlocked
}

func (a User) GetPasswordChangedAt() *time.Time {
	return a.PasswordChangedAt
}

func (a *User) GenerateToken(db *gorm.DB) (string, time.Time, error) {
	expStr := config.Env.JWTExpiration
	days, err := strconv.Atoi(strings.TrimSuffix(expStr, "d"))
	if err != nil {
		return "", time.Time{}, err
	}
	now := time.Now()
	tokenExpDate := now.AddDate(0, 0, days)

	claims := jwt.MapClaims{
		"userId": a.ID.String(),
		"role":   a.Role,
		"iat":    now.Unix(),
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

func (a *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}
