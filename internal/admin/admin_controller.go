package admin

import (
	"errors"
	"net/http"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMyProfile(c *gin.Context, db *gorm.DB) {
	var admin Admin

	adminId := utils.GetUserID(c)

	result := db.First(&admin, "id=?", adminId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Profile not found", http.StatusNotFound))
		} else {
			c.Error(result.Error)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile fetched successfully",
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
