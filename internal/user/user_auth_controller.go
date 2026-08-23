package user

import (
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
