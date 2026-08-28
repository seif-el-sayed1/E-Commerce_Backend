package subcategory

import (
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateSubCategoryRequest struct {
	Name       string    `json:"name" binding:"required"`
	CategoryID uuid.UUID `json:"category_id" binding:"required"`
}

type UpdateSubCategoryRequest struct {
	Name       *string    `json:"name"`
	CategoryID *uuid.UUID `json:"category_id"`
}

func ValidateCreateSubCategory(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

func ValidateUpdateSubCategory(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}
