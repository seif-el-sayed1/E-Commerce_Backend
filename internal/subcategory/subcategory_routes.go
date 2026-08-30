package subcategory

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"
)

func SubCategoryRoutes(rg *gin.RouterGroup, db *gorm.DB, protect gin.HandlerFunc, allowedTo func(roles ...string) gin.HandlerFunc) {
	subcats := rg.Group("/subcategories")
	{
		// Public
		subcats.GET("", func(c *gin.Context) {
			GetSubCategories(c, db)
		})
		subcats.GET("/:id", func(c *gin.Context) {
			GetSubCategory(c, db)
		})

		subcats.POST("", protect, allowedTo(constants.Roles.Admin, constants.Roles.SuperAdmin), func(c *gin.Context) {
			CreateSubCategory(c, db)
		})
		subcats.PUT("/:id", protect, allowedTo(constants.Roles.Admin, constants.Roles.SuperAdmin), func(c *gin.Context) {
			UpdateSubCategory(c, db)
		})
		subcats.DELETE("/:id", protect, allowedTo(constants.Roles.Admin, constants.Roles.SuperAdmin), func(c *gin.Context) {
			DeleteSubCategory(c, db)
		})
	}
}
