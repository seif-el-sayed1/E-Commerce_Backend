package admin

import (
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func ValidateAdminLogin(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
}

func ValidateAdminChangePassword(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type VerifyAccount struct {
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
}

func ValidateAdminVerifyAccount(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type ForgetPassword struct {
	Email string `json:"email" binding:"required,email"`
}

func ValidateAdminForgetPassword(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type ResetPassword struct {
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
}

func ValidateAdminResetPassword(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}
