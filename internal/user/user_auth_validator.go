package user

import (
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
}

func ValidateUserLogin(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}

	req, ok := obj.(*UserLoginRequest)
	if ok && req.Email == "" && req.Phone == "" {
		c.Error(utils.NewApiError("Email or phone is required", http.StatusBadRequest))
		c.Abort()
		return false
	}

	return true
}

type UserRegisterRequest struct {
	FirstName string  `json:"firstName" binding:"required"`
	LastName  string  `json:"lastName"  binding:"required"`
	Email     string  `json:"email"     binding:"required,email"`
	Password  string  `json:"password"  binding:"required,min=8"`
	Phone     *string `json:"phone"`
}

func ValidateUserRegister(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type UserVerifyAccountRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code"  binding:"required"`
}

func ValidateUserVerifyAccount(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

func ValidateUpdatePassword(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type VerifyOTPRequest struct {
	OTP string `json:"otp" binding:"required"`
}

func ValidateVerifyOTP(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func ValidateSendOTP(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}
