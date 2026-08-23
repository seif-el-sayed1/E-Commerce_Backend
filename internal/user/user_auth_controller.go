package user

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/email"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func login(c *gin.Context, db *gorm.DB, user User, password string) {
	if !user.ComparePassword(password) {
		c.Error(utils.NewApiError("Incorrect email or password", http.StatusUnauthorized))
		c.Abort()
		return
	}

	message := "Welcome back " + user.FirstName + "!"

	if !user.IsActive {
		if err := db.Model(&user).Updates(map[string]interface{}{
			"is_active": true,
		}).Error; err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		message = "Welcome back! Your account has been reactivated."
	}

	if !user.IsVerified {
		rawCode, hashedCode, err := utils.GenerateOTPCode()
		if err != nil {
			c.Error(utils.NewApiError("Failed to generate verification code", http.StatusInternalServerError))
			c.Abort()
			return
		}

		exp := time.Now().Add(10 * time.Minute)
		if err := db.Model(&user).Updates(map[string]interface{}{
			"verification_code":     hashedCode,
			"verification_code_exp": exp,
		}).Error; err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if err := email.UserVerificationEmail(rawCode, user.Email); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Verification OTP is sent to your Email",
			"data": gin.H{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
				"phone":      user.Phone,
				"role":       user.Role,
			},
		})
		return
	}

	if user.IsBlocked {
		c.Error(utils.NewApiError("Your account is blocked, please contact the support team", http.StatusForbidden))
		c.Abort()
		return
	}

	token, expDate, err := user.GenerateToken(db)
	if err != nil {
		c.Error(utils.NewApiError("Failed to generate token", http.StatusInternalServerError))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"token":   token,
		"exp":     expDate,
		"data": gin.H{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"phone":      user.Phone,
			"role":       user.Role,
		},
	})
}

// @desc    Log In
// @route   POST /user/auth/login
// @access  Public
func UserLogin(c *gin.Context, db *gorm.DB) {
	var body UserLoginRequest
	if !ValidateUserLogin(c, &body) {
		return
	}

	var user User
	var result *gorm.DB
	if body.Phone != "" {
		result = db.First(&user, "phone = ?", body.Phone)
	} else {
		result = db.First(&user, "email = ?", body.Email)
	}

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Incorrect email or password", http.StatusUnauthorized))
		} else {
			c.Error(result.Error)
		}
		c.Abort()
		return
	}

	login(c, db, user, body.Password)
}

// @desc    Sign Up
// @route   POST /user/auth/register
// @access  Public
func UserRegister(c *gin.Context, db *gorm.DB) {
	var body UserRegisterRequest
	if !ValidateUserRegister(c, &body) {
		return
	}

	rawCode, hashedCode, err := utils.GenerateOTPCode()
	if err != nil {
		c.Error(utils.NewApiError("Failed to generate verification code", http.StatusInternalServerError))
		c.Abort()
		return
	}

	exp := time.Now().Add(10 * time.Minute)

	newUser := User{
		FirstName:           body.FirstName,
		LastName:            body.LastName,
		Email:               body.Email,
		Phone:               body.Phone,
		Password:            body.Password,
		VerificationCode:    &hashedCode,
		VerificationCodeExp: &exp,
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if err := email.UserVerificationEmail(rawCode, newUser.Email); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification OTP is sent to your Email",
		"data": gin.H{
			"id":         newUser.ID,
			"first_name": newUser.FirstName,
			"last_name":  newUser.LastName,
			"email":      newUser.Email,
			"phone":      newUser.Phone,
			"role":       newUser.Role,
		},
	})
}

// @desc    User account verification
// @route   POST /user/auth/verify-account
// @access  Public
func UserVerifyAccount(c *gin.Context, db *gorm.DB) {
	var body UserVerifyAccountRequest
	if !ValidateUserVerifyAccount(c, &body) {
		return
	}

	// Find user by email
	var user User
	if err := db.First(&user, "email = ?", body.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Invalid request", http.StatusBadRequest))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	// Make sure a verification code was actually issued
	if user.VerificationCode == nil || user.VerificationCodeExp == nil {
		c.Error(utils.NewApiError("Invalid request", http.StatusBadRequest))
		c.Abort()
		return
	}

	// Check expiry
	if time.Now().After(*user.VerificationCodeExp) {
		c.Error(utils.NewApiError("Verification OTP has expired", http.StatusUnauthorized))
		c.Abort()
		return
	}

	// Hash the incoming code and compare
	hash := sha256.Sum256([]byte(body.Code))
	hashedCode := hex.EncodeToString(hash[:])
	if *user.VerificationCode != hashedCode {
		c.Error(utils.NewApiError("Invalid verification OTP", http.StatusUnauthorized))
		c.Abort()
		return
	}

	// Mark account as verified and clear the OTP fields
	if err := db.Model(&user).Updates(map[string]interface{}{
		"is_verified":           true,
		"verification_code":     nil,
		"verification_code_exp": nil,
	}).Error; err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	// Generate JWT token
	token, expDate, err := user.GenerateToken(db)
	if err != nil {
		c.Error(utils.NewApiError("Failed to generate token", http.StatusInternalServerError))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account verified successfully",
		"token":   token,
		"exp":     expDate,
		"data": gin.H{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"phone":      user.Phone,
			"role":       user.Role,
		},
	})
}

// @desc    Verify OTP (used when user has a code but no token yet)
// @route   POST /user/auth/verify-otp
// @access  Public
func VerifyOTP(c *gin.Context, db *gorm.DB) {
	var body VerifyOTPRequest
	if !ValidateVerifyOTP(c, &body) {
		return
	}

	// Hash the incoming OTP and look for a matching, non-expired record
	hash := sha256.Sum256([]byte(body.OTP))
	hashedCode := hex.EncodeToString(hash[:])

	var user User
	err := db.Where(
		"verification_code = ? AND verification_code_exp > ?",
		hashedCode, time.Now(),
	).First(&user).Error
	if err != nil {
		c.Error(utils.NewApiError("OTP isn't found or has expired", http.StatusForbidden))
		c.Abort()
		return
	}

	// Mark as verified and clear OTP fields
	if err := db.Model(&user).Updates(map[string]interface{}{
		"is_verified":           true,
		"verification_code":     nil,
		"verification_code_exp": nil,
	}).Error; err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	// Reuse existing token if present, otherwise generate a new one
	var token string
	var expDate time.Time
	if user.Token != nil {
		token = *user.Token
		if user.TokenExpDate != nil {
			expDate = *user.TokenExpDate
		}
	} else {
		token, expDate, err = user.GenerateToken(db)
		if err != nil {
			c.Error(utils.NewApiError("Failed to generate token", http.StatusInternalServerError))
			c.Abort()
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account verified successfully",
		"token":   token,
		"exp":     expDate,
		"data": gin.H{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"phone":      user.Phone,
			"role":       user.Role,
		},
	})
}

// @desc    Send OTP to an unverified user's email
// @route   POST /user/auth/send-otp
// @access  Public
func SendOTP(c *gin.Context, db *gorm.DB) {
	var body SendOTPRequest
	if !ValidateSendOTP(c, &body) {
		return
	}

	// Find unverified user by email
	var user User
	err := db.Where("email = ? AND is_verified = ?", body.Email, false).First(&user).Error
	if err != nil {
		c.Error(utils.NewApiError("User not found", http.StatusNotFound))
		c.Abort()
		return
	}

	// Generate new OTP
	rawCode, hashedCode, err := utils.GenerateOTPCode()
	if err != nil {
		c.Error(utils.NewApiError("Failed to generate OTP", http.StatusInternalServerError))
		c.Abort()
		return
	}

	exp := time.Now().Add(10 * time.Minute)
	if err := db.Model(&user).Updates(map[string]interface{}{
		"verification_code":     hashedCode,
		"verification_code_exp": exp,
	}).Error; err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if err := email.UserVerificationEmail(rawCode, user.Email); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification OTP is sent to your Email",
	})
}

// @desc    Update logged user password
// @route   PATCH /user/auth/update-password
// @access  Private
func UpdateLoggedUserPassword(c *gin.Context, db *gorm.DB) {
	currentUser := utils.GetUserData(c).(User)

	var body UpdatePasswordRequest
	if !ValidateUpdatePassword(c, &body) {
		return
	}

	// Verify current password
	if !currentUser.ComparePassword(body.CurrentPassword) {
		c.Error(utils.NewApiError("Incorrect password", http.StatusUnauthorized))
		c.Abort()
		return
	}

	// Update with new password — BeforeUpdate hook hashes it automatically
	currentUser.Password = body.NewPassword

	if err := db.Save(&currentUser).Error; err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password updated successfully, please login again",
	})
}

// @desc    Log out — clears the stored token
// @route   POST /user/auth/log-out
// @access  Private
func UserLogout(c *gin.Context, db *gorm.DB) {
	currentUser := utils.GetUserData(c).(User)

	if err := db.Model(&currentUser).Updates(map[string]interface{}{
		"token":          nil,
		"token_exp_date": nil,
	}).Error; err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User logged out successfully!",
	})
}
