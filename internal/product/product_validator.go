package product

import (
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name          string    `json:"name" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	SubCategoryID uuid.UUID `json:"subcategory_id" binding:"required"`
}

type UpdateProductRequest struct {
	Name          *string    `json:"name"`
	Description   *string    `json:"description"`
	SubCategoryID *uuid.UUID `json:"subcategory_id"`
}

func ValidateCreateProduct(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

func ValidateUpdateProduct(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}
