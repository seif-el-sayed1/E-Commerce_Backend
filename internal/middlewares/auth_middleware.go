package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/admin"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"
)

// Claims
type Claims struct {
	UserID uuid.UUID `json:"userId"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// AuthUser defines the common authentication requirements for all user types.
type AuthUser interface {
	GetToken() *string
	GetIsActive() bool
	GetIsBlocked() bool
	GetPasswordChangedAt() *time.Time
}

func checkUser[T AuthUser](db *gorm.DB, token string, claims *Claims) (T, *utils.ApiError) {
	var currentUser T

	if err := db.First(&currentUser, "id = ?", claims.UserID).Error; err != nil {
		return currentUser, utils.NewApiError("not found", 401)
	}

	// Check if the token matches the one stored in the database
	if currentUser.GetToken() == nil || *currentUser.GetToken() != token {
		return currentUser, utils.NewApiError("Session expired, please login again...", 401)
	}

	// Check if the user is active
	if !currentUser.GetIsActive() {
		return currentUser, utils.NewApiError("account is deactivated", 401)
	}

	// Check if the user is blocked
	if currentUser.GetIsBlocked() {
		return currentUser, utils.NewApiError("Your account is blocked, please contact the support team", 401)
	}

	// Check if the password has been changed
	if pca := currentUser.GetPasswordChangedAt(); pca != nil {
		if pca.Unix() > claims.IssuedAt.Unix() {
			return currentUser, utils.NewApiError("Password recently changed, please login again...", 401)
		}
	}

	return currentUser, nil
}

func Protect(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")
		if token == "" {
			abortWithError(c, utils.NewApiError("Invalid token, please login again...", 401))
			return
		}

		claims := &Claims{}

		// Parse and validate the token
		parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Env.JWTSecret), nil
		})
		if err != nil || !parsed.Valid {
			abortWithError(c, utils.NewApiError("Invalid token, please login again...", 401))
			return
		}

		// Check if the role is not in app roles
		if !contains(constants.RolesList, claims.Role) {
			abortWithError(c, utils.NewApiError("Invalid token role, please login again...", 401))
			return
		}

		// Check if the token has expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			abortWithError(c, utils.NewApiError("Token has expired, please login again...", 401))
			return
		}

		var currentUser interface{}
		var apiErr *utils.ApiError

		// Check the roles
		switch claims.Role {
		// ADMIN, SUPER_ADMIN
		case constants.Roles.SuperAdmin, constants.Roles.Admin:
			currentUser, apiErr = checkUser[admin.Admin](db, token, claims)

		default:
			apiErr = utils.NewApiError("Invalid token role, please login again...", 401)
		}

		if apiErr != nil {
			abortWithError(c, apiErr)
			return
		}

		c.Set("role", claims.Role)
		c.Set("userId", claims.UserID)
		c.Set("user", currentUser)
		c.Next()
	}
}

// helpers
func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func abortWithError(c *gin.Context, err *utils.ApiError) {
	c.AbortWithStatusJSON(err.StatusCode, gin.H{"message": err.Message})
}
