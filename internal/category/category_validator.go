package category

import (
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateCategoryRequest struct {
	Name    string `json:"name" binding:"required"`
	CatType string `json:"category_type" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name    *string `json:"name"`
	CatType *string `json:"category_type"`
}

func ValidateCreateCategory(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

func ValidateUpdateCategory(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}
