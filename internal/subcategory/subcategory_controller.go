package subcategory

import (
	"errors"
	"net/http"
	"strings"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/category"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @desc    Create subcategory
// @route   POST /subcategories
// @access  Private (Admin)
func CreateSubCategory(c *gin.Context, db *gorm.DB) {
	var body CreateSubCategoryRequest
	if !ValidateCreateSubCategory(c, &body) {
		return
	}

	// Verify category exists
	var cat category.Category
	if err := db.First(&cat, "id = ?", body.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Category not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	subCat := SubCategory{
		Name:       body.Name,
		CategoryID: body.CategoryID,
	}

	if err := db.Create(&subCat).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value") {
			c.Error(utils.NewApiError("Subcategory already exists", http.StatusBadRequest))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    subCat,
	})
}

// @desc    Get all subcategories
// @route   GET /subcategories
// @access  Public
func GetSubCategories(c *gin.Context, db *gorm.DB) {
	query := db.Model(&SubCategory{}).Preload("Category")

	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}

	features := utils.New(c, query, "SubCategory").
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

	var subCategories []SubCategory
	if err := features.Execute(&subCategories); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"pagination": features.Pagination,
		"data":       subCategories,
	})
}

// @desc    Get specific subcategory
// @route   GET /subcategories/:id
// @access  Public
func GetSubCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var subCat SubCategory
	if err := db.Preload("Category").First(&subCat, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Subcategory not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subCat,
	})
}

// @desc    Update specific subcategory
// @route   PUT /subcategories/:id
// @access  Private (Admin)
func UpdateSubCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var body UpdateSubCategoryRequest
	if !ValidateUpdateSubCategory(c, &body) {
		return
	}

	var subCat SubCategory
	if err := db.First(&subCat, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Subcategory not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	if body.CategoryID != nil {
		var cat category.Category
		if err := db.First(&cat, "id = ?", *body.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Error(utils.NewApiError("Category not found", http.StatusNotFound))
			} else {
				c.Error(err)
			}
			c.Abort()
			return
		}
		subCat.CategoryID = *body.CategoryID
	}

	if body.Name != nil {
		subCat.Name = *body.Name
	}

	if err := db.Save(&subCat).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value") {
			c.Error(utils.NewApiError("Subcategory name already exists", http.StatusBadRequest))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	db.Preload("Category").First(&subCat, "id = ?", subCat.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subCat,
	})
}

// @desc    Delete specific subcategory
// @route   DELETE /subcategories/:id
// @access  Private (Admin)
func DeleteSubCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var subCat SubCategory
	if err := db.First(&subCat, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Subcategory not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	if err := db.Delete(&subCat).Error; err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Subcategory deleted successfully",
	})
}
