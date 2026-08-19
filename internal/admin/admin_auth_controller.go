package admin

import (
	"errors"
	"net/http"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/email"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @desc    Admin login
// @route   POST /admins/auth/login
// @access  Public
func Login(c *gin.Context, db *gorm.DB) {
	var adminLoginValidator AdminLogin

	if !ValidateAdminLogin(c, &adminLoginValidator) {
		return
	}

	var admin Admin
	result := db.First(&admin, "email = ?", adminLoginValidator.Email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Incorrect email or password", http.StatusUnauthorized))
		} else {
			c.Error(result.Error)
		}
		c.Abort()
		return
	}

	if admin.IsDeleted {
		c.Error(utils.NewApiError("Incorrect email or password", http.StatusForbidden))
		c.Abort()
		return
	}

	if !admin.ComparePassword(adminLoginValidator.Password) {
		c.Error(utils.NewApiError("Incorrect email or password", http.StatusUnauthorized))
		c.Abort()
		return
	}

	token, expDate, err := admin.GenerateToken(db)
	if err != nil {
		c.Error(utils.NewApiError("Failed to generate token", http.StatusInternalServerError))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
		"exp":     expDate,
		"data": gin.H{
			"id":             admin.ID,
			"first_name":     admin.FirstName,
			"last_name":      admin.LastName,
			"email":          admin.Email,
			"is_super_admin": admin.IsSuperAdmin,
			"role":           admin.Role,
			"is_active":      admin.IsActive,
			"is_verified":    admin.IsVerified,
		},
	})

}

// @desc    Change logged in admin password
// @route   PATCH /adminS/auth/change-password
// @access  Private
func AdminChangePassword(c *gin.Context, db *gorm.DB) {
	currentAdmin := utils.GetUserData(c).(Admin)

	var body ChangePassword
	if !ValidateAdminChangePassword(c, &body) {
		return
	}

	if !currentAdmin.ComparePassword(body.CurrentPassword) {
		c.Error(utils.NewApiError("Incorrect password", http.StatusBadRequest))
		c.Abort()
		return
	}

	if body.NewPassword != body.ConfirmPassword {
		c.Error(utils.NewApiError("Passwords do not match", http.StatusBadRequest))
		c.Abort()
		return
	}

	currentAdmin.Password = body.NewPassword

	now := time.Now()
	currentAdmin.PasswordChangedAt = &now

	result := db.Save(&currentAdmin)
	if result.Error != nil {
		c.Error(result.Error)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully",
	})

}

// @desc    Verify logged in admin account
// @route   POST /adminS/auth/verify-account
// @access  Private
func VerifyAdminAccount(c *gin.Context, db *gorm.DB) {
	currentAdmin := utils.GetUserData(c).(Admin)

	var body VerifyAccount
	if !ValidateAdminVerifyAccount(c, &body) {
		return
	}

	if currentAdmin.IsVerified {
		c.Error(utils.NewApiError("Account already verified", http.StatusBadRequest))
		c.Abort()
		return
	}

	if body.Password != body.ConfirmPassword {
		c.Error(utils.NewApiError("Passwords do not match", http.StatusBadRequest))
		c.Abort()
		return
	}

	currentAdmin.IsVerified = true
	currentAdmin.Password = body.Password

	now := time.Now()
	currentAdmin.PasswordChangedAt = &now

	result := db.Save(&currentAdmin)
	if result.Error != nil {
		c.Error(result.Error)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account verified successfully",
	})

}

// @desc    Forgot admin account password
// @route   POST /admins/auth/forgot-Password
// @access  Public
func AdminForgetPassword(c *gin.Context, db *gorm.DB) {
	var body ForgetPassword
	if !ValidateAdminForgetPassword(c, &body) {
		return
	}

	var admin Admin
	result := db.First(&admin, "email = ?", body.Email)
	if result.Error != nil {
		c.Error(result.Error)
		c.Abort()
		return
	}

	if admin.IsBlocked || admin.IsDeleted {
		c.Error(utils.NewApiError("Account is blocked or deleted", http.StatusForbidden))
		c.Abort()
		return
	}

	token := admin.CreatePasswordResetToken()

	result = db.Save(&admin)
	if result.Error != nil {
		c.Error(result.Error)
		c.Abort()
		return
	}

	err := email.AdminForgotPasswordEmail(token, body.Email)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset email sent successfully",
	})
}

// @desc    Forgot admin account password
// @route   PATCH /admin/auth/reset-password/:token
// @access  Private
func AdminResetPassword(c *gin.Context, db *gorm.DB) {
	var body ResetPassword
	token := c.Param("token")

	if !ValidateAdminResetPassword(c, &body) {
		return
	}

	if body.Password != body.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Passwords do not match",
		})
		return
	}

	hashedToken := sha256Hex(token)

	var admin Admin
	result := db.First(
		&admin,
		"password_reset_token = ? AND password_reset_expires_at > ?",
		hashedToken, time.Now(),
	)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Token is invalid or has expired",
		})
		return
	}

	now := time.Now()
	admin.Password = body.Password
	admin.PasswordChangedAt = &now
	admin.VerificationToken = nil
	admin.PasswordResetToken = nil
	admin.PasswordResetExpiresAt = nil

	result = db.Save(&admin)
	if result.Error != nil {
		c.Error(result.Error)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset successfully",
	})
}
