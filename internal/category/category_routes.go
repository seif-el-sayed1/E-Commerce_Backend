package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"
)

func CategoryRoutes(rg *gin.RouterGroup, db *gorm.DB, protect gin.HandlerFunc, allowedTo func(roles ...string) gin.HandlerFunc) {
	categories := rg.Group("/categories")
	{
		// Public
		categories.GET("", func(c *gin.Context) {
			GetCategories(c, db)
		})
		categories.GET("/:id", func(c *gin.Context) {
			GetCategory(c, db)
		})

		categories.POST("", protect, allowedTo(constants.Roles.Admin, constants.Roles.SuperAdmin), func(c *gin.Context) {
			CreateCategory(c, db)
		})
		categories.PUT("/:id", protect, allowedTo(constants.Roles.Admin, constants.Roles.SuperAdmin), func(c *gin.Context) {
			UpdateCategory(c, db)
		})
		categories.DELETE("/:id", protect, allowedTo(constants.Roles.Admin, constants.Roles.SuperAdmin), func(c *gin.Context) {
			DeleteCategory(c, db)
		})
	}
}
