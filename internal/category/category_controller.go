package category

import (
	"errors"
	"net/http"
	"strings"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @desc    Create category
// @route   POST /categories
// @access  Private (Admin)
func CreateCategory(c *gin.Context, db *gorm.DB) {
	var body CreateCategoryRequest
	if !ValidateCreateCategory(c, &body) {
		return
	}

	cat := Category{
		Name: body.Name,
	}
	if body.CatType != "" {
		cat.CatType = body.CatType
	}

	if err := db.Create(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value") {
			c.Error(utils.NewApiError("Category already exists", http.StatusBadRequest))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    cat,
	})
}

// @desc    Get all categories
// @route   GET /categories
// @access  Public
func GetCategories(c *gin.Context, db *gorm.DB) {
	features := utils.New(c, db.Model(&Category{}), "Category").
		Search().
		Filter().
		Sort().
		Paginate()

	features, err := features.CalculatePagination()
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var categories []Category
	if err := features.Execute(&categories); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"pagination": features.Pagination,
		"data":       categories,
	})
}

// @desc    Get specific category
// @route   GET /categories/:id
// @access  Public
func GetCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var cat Category
	if err := db.First(&cat, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Category not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cat,
	})
}

// @desc    Update specific category
// @route   PUT /categories/:id
// @access  Private (Admin)
func UpdateCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var body UpdateCategoryRequest
	if !ValidateUpdateCategory(c, &body) {
		return
	}

	var cat Category
	if err := db.First(&cat, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Category not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	if body.Name != nil {
		cat.Name = *body.Name
	}
	if body.CatType != nil {
		cat.CatType = *body.CatType
	}

	if err := db.Save(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value") {
			c.Error(utils.NewApiError("Category name already exists", http.StatusBadRequest))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cat,
	})
}

// @desc    Delete specific category
// @route   DELETE /categories/:id
// @access  Private (Admin)
func DeleteCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var cat Category
	if err := db.First(&cat, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Category not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	if err := db.Delete(&cat).Error; err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category deleted successfully",
	})
}
