package admin

import (
	"errors"
	"net/http"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Add struct {
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Phone     *string `json:"phone" binding:"required"`
}

func ValidateAddAdmin(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			message := validationMessage(ve[0])
			c.Error(utils.NewApiError(message, http.StatusBadRequest))
			c.Abort()
			return false
		}

		c.Error(utils.NewApiError("Invalid request body", http.StatusBadRequest))
		c.Abort()
		return false
	}
	return true
}
