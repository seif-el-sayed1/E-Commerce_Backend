package admin

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (a Admin) GetToken() *string {
	return a.Token
}

func (a Admin) GetIsActive() bool {
	return a.IsActive
}

func (a Admin) GetIsBlocked() bool {
	return a.IsBlocked
}

func (a Admin) GetPasswordChangedAt() *time.Time {
	return a.PasswordChangedAt
}

func (a *Admin) GenerateToken(db *gorm.DB) (string, time.Time, error) {
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
func (a *Admin) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashed)

	changedAt := time.Now().Add(-1 * time.Second)
	a.PasswordChangedAt = &changedAt

	return nil
}

func (a *Admin) BeforeUpdate(tx *gorm.DB) error {
	if len(a.Password) > 0 && !strings.HasPrefix(a.Password, "$2a$") &&
		!strings.HasPrefix(a.Password, "$2b$") && !strings.HasPrefix(a.Password, "$2y$") {

		hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		a.Password = string(hashed)

		changedAt := time.Now().Add(-1 * time.Second)
		a.PasswordChangedAt = &changedAt
	}
	return nil
}
func (a *Admin) CreateEmailToken(forPassword bool) (rawToken string, resetToken *string) {
	raw := generateRandomHex(32)

	hashed := sha256Hex(raw)
	a.VerificationToken = &hashed

	expiresAt := time.Now().Add(10 * time.Minute)
	a.VerificationTokenExp = &expiresAt

	if forPassword {
		return raw, &hashed
	}
	return raw, nil
}

func (a *Admin) CreatePasswordResetToken() string {
	rawToken, resetToken := a.CreateEmailToken(true)

	a.PasswordResetToken = resetToken
	expiresAt := time.Now().Add(10 * time.Minute)
	a.PasswordResetExpiresAt = &expiresAt

	return rawToken
}

// helpers
func generateRandomHex(n int) string {
	bytes := make([]byte, n)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func sha256Hex(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
