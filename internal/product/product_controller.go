package product

import (
	"errors"
	"net/http"
	"strings"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/subcategory"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @desc    Create product
// @route   POST /products
// @access  Private (Admin)
func CreateProduct(c *gin.Context, db *gorm.DB) {
	var body CreateProductRequest
	if !ValidateCreateProduct(c, &body) {
		return
	}

	// Verify subcategory exists (with its category)
	var subCat subcategory.SubCategory
	if err := db.Preload("Category").First(&subCat, "id = ?", body.SubCategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("SubCategory not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	prod := Product{
		Name:          body.Name,
		Description:   body.Description,
		SubCategoryID: body.SubCategoryID,
		SubCategory:   subCat,
	}

	if err := db.Omit("SubCategory").Create(&prod).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value") {
			c.Error(utils.NewApiError("Product name already exists", http.StatusBadRequest))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    prod,
	})
}

// @desc    Get all products
// @route   GET /products
// @access  Public
func GetProducts(c *gin.Context, db *gorm.DB) {
	query := db.Model(&Product{}).Preload("SubCategory.Category")

	// Support filtering by subcategory ID
	if subCatID := c.Query("subcategory_id"); subCatID != "" {
		query = query.Where("subcategory_id = ?", subCatID)
	}

	features := utils.New(c, query, "Product").
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

	var products []Product
	if err := features.Execute(&products); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"pagination": features.Pagination,
		"data":       products,
	})
}

// @desc    Get specific product
// @route   GET /products/:id
// @access  Public
func GetProduct(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var prod Product
	if err := db.Preload("SubCategory.Category").First(&prod, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Product not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prod,
	})
}

// @desc    Update specific product
// @route   PUT /products/:id
// @access  Private (Admin)
func UpdateProduct(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var body UpdateProductRequest
	if !ValidateUpdateProduct(c, &body) {
		return
	}

	var prod Product
	if err := db.First(&prod, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Product not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	if body.SubCategoryID != nil {
		// Verify subcategory exists
		var subCat subcategory.SubCategory
		if err := db.First(&subCat, "id = ?", *body.SubCategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Error(utils.NewApiError("SubCategory not found", http.StatusNotFound))
			} else {
				c.Error(err)
			}
			c.Abort()
			return
		}
		prod.SubCategoryID = *body.SubCategoryID
	}

	if body.Name != nil {
		prod.Name = *body.Name
	}
	if body.Description != nil {
		prod.Description = *body.Description
	}

	if err := db.Omit("SubCategory").Save(&prod).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value") {
			c.Error(utils.NewApiError("Product name already exists", http.StatusBadRequest))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	db.Preload("SubCategory.Category").First(&prod, "id = ?", prod.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prod,
	})
}

// @desc    Delete specific product
// @route   DELETE /products/:id
// @access  Private (Admin)
func DeleteProduct(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var prod Product
	if err := db.First(&prod, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NewApiError("Product not found", http.StatusNotFound))
		} else {
			c.Error(err)
		}
		c.Abort()
		return
	}

	if err := db.Delete(&prod).Error; err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product deleted successfully",
	})
}
