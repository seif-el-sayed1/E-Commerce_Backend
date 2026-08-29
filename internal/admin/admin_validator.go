package admin

import (
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
)

type Add struct {
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Phone     *string `json:"phone" binding:"required"`
}

func ValidateAddAdmin(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}

type Update struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Phone     *string `json:"phone"`
}

func ValidateUpdateAdmin(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(utils.NewApiError(utils.FormatValidationError(err), http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}
