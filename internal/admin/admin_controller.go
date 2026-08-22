package admin

import (
	"errors"
	"net/http"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/email"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var adminSelectColumns = "id, first_name, last_name, email, phone, role, is_super_admin, is_verified, is_blocked, is_deleted, is_active, created_at"

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

// @desc    Get all admins
// @route   GET /admins
// @access  Private
func GetAllAdmins(c *gin.Context, db *gorm.DB) {
	var admins []Admin

	af := utils.New(c, db.Model(&Admin{}).Select(adminSelectColumns), "Admin")
	af.Search().Filter().Sort().Paginate()

	af, err := af.CalculatePagination()
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if err := af.Execute(&admins); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	data := make([]gin.H, 0, len(admins))
	for _, admin := range admins {
		data = append(data, gin.H{
			"id":             admin.ID,
			"first_name":     admin.FirstName,
			"last_name":      admin.LastName,
			"email":          admin.Email,
			"phone":          admin.Phone,
			"role":           admin.Role,
			"is_super_admin": admin.IsSuperAdmin,
			"is_verified":    admin.IsVerified,
			"is_blocked":     admin.IsBlocked,
			"is_deleted":     admin.IsDeleted,
			"is_active":      admin.IsActive,
			"created_at":     admin.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Admins fetched successfully",
		"data":       data,
		"pagination": af.Pagination,
	})
}

// @desc    Get single admin
// @route   GET /admins/:id
// @access  Private
func GetSingleAdmin(c *gin.Context, db *gorm.DB) {
	adminId := c.Param("id")
	var admin Admin

	result := db.Select(adminSelectColumns).First(&admin, "id=?", adminId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Admin not found", http.StatusNotFound))
		} else {
			c.Error(result.Error)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin fetched successfully",
		"data": gin.H{
			"id":             admin.ID,
			"first_name":     admin.FirstName,
			"last_name":      admin.LastName,
			"email":          admin.Email,
			"phone":          admin.Phone,
			"role":           admin.Role,
			"is_super_admin": admin.IsSuperAdmin,
			"is_verified":    admin.IsVerified,
			"is_blocked":     admin.IsBlocked,
			"is_deleted":     admin.IsDeleted,
			"is_active":      admin.IsActive,
			"created_at":     admin.CreatedAt,
		},
	})
}

// @desc    Get my profile
// @route   GET /admins/profile
// @access  Private
func GetMyProfile(c *gin.Context, db *gorm.DB) {
	var admin Admin

	adminId := utils.GetUserID(c)

	result := db.Select(adminSelectColumns).First(&admin, "id=?", adminId)
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
			"role":           admin.Role,
			"is_super_admin": admin.IsSuperAdmin,
			"is_verified":    admin.IsVerified,
			"is_blocked":     admin.IsBlocked,
			"is_deleted":     admin.IsDeleted,
			"is_active":      admin.IsActive,
			"created_at":     admin.CreatedAt,
		},
	})
}
