package user

import (
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @desc    Get my profile
// @route   GET /users/me
// @access  Private (User)
func GetMyProfile(c *gin.Context, db *gorm.DB) {
	// The protect middleware sets the "user" in the context
	currentUser := utils.GetUserData(c).(User)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":         currentUser.ID,
			"first_name": currentUser.FirstName,
			"last_name":  currentUser.LastName,
			"email":      currentUser.Email,
			"phone":      currentUser.Phone,
			"role":       currentUser.Role,
		},
	})
}

// @desc    Update my profile
// @route   PATCH /users/me
// @access  Private (User)
func UpdateMyProfile(c *gin.Context, db *gorm.DB) {
	currentUser := utils.GetUserData(c).(User)

	var body UpdateProfileRequest
	if !ValidateUpdateProfile(c, &body) {
		return
	}

	// Update only allowed fields
	if body.FirstName != nil {
		currentUser.FirstName = *body.FirstName
	}
	if body.LastName != nil {
		currentUser.LastName = *body.LastName
	}
	if body.Phone != nil {
		currentUser.Phone = body.Phone
	}

	if err := db.Save(&currentUser).Error; err != nil {
		c.Error(utils.NewApiError("Failed to update profile", http.StatusInternalServerError))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":         currentUser.ID,
			"first_name": currentUser.FirstName,
			"last_name":  currentUser.LastName,
			"email":      currentUser.Email,
			"phone":      currentUser.Phone,
			"role":       currentUser.Role,
		},
	})
}
