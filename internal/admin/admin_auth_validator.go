package admin

import (
	"errors"
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AdminLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func ValidateAdminLogin(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			message := validationMessage(ve[0])
			c.Error(utils.NewApiError(message, http.StatusBadRequest))
			c.Abort()
			return false
		}

		c.Error(utils.NewApiError("Invalid request body", http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type AdminChangePasword struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
}

func ValidateAdminChangePassword(c *gin.Context, obj interface{}) bool {

	if err := c.ShouldBindJSON(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			message := validationMessage(ve[0])
			c.Error(utils.NewApiError(message, http.StatusBadRequest))
			c.Abort()
			return false
		}

		c.Error(utils.NewApiError("Invalid request body", http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return utils.CamelToWords(fe.Field()) + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return utils.CamelToWords(fe.Field()) + " must have at least " + fe.Param() + " characters"
	default:
		return "Invalid value for " + utils.CamelToWords(fe.Field())
	}
}
