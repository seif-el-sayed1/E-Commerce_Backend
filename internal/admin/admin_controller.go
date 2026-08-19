package admin

import (
	"net/http"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/email"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @desc    Add new admin
// @route   POST /admins/add-admin
// @access  Private
func AddAdmin(c *gin.Context, db *gorm.DB) {
	var body Add
	if !ValidateAddAdmin(c, &body) {
		return
	}
	var admin Admin
	admin.FirstName = body.FirstName
	admin.LastName = body.LastName
	admin.Email = body.Email
	admin.Phone = body.Phone

	result := db.Create(&admin)
	if result.Error != nil {
		c.Error(result.Error)
		c.Abort()
		return
	}

	token, _, err := admin.GenerateToken(db)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	err = email.AdminVerificationEmail(token, body.Email)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin added successfully",
		"data": gin.H{
			"id":             admin.ID,
			"first_name":     admin.FirstName,
			"last_name":      admin.LastName,
			"email":          admin.Email,
			"phone":          admin.Phone,
			"is_super_admin": admin.IsSuperAdmin,
			"role":           admin.Role,
			"is_active":      admin.IsActive,
			"is_verified":    admin.IsVerified,
		},
	})

}
