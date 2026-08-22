package user

import (
	"errors"
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UpdateProfileRequest struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
}

func ValidateUpdateProfile(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.Error(utils.NewApiError(userValidationMessage(ve[0]), http.StatusBadRequest))
			c.Abort()
			return false
		}
		c.Error(utils.NewApiError("Invalid request body", http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}
