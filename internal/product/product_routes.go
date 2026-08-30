package product

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"
)

func ProductRoutes(rg *gin.RouterGroup, db *gorm.DB, protect gin.HandlerFunc, allowedTo func(roles ...string) gin.HandlerFunc) {
	products := rg.Group("/products")
	{
		// Public
		products.GET("", func(c *gin.Context) {
			GetProducts(c, db)
		})
		products.GET("/:id", func(c *gin.Context) {
			GetProduct(c, db)
		})

		products.POST("", protect, allowedTo(constants.Roles.Admin, constants.Roles.SuperAdmin), func(c *gin.Context) {
			CreateProduct(c, db)
		})
		products.PUT("/:id", protect, allowedTo(constants.Roles.Admin, constants.Roles.SuperAdmin), func(c *gin.Context) {
			UpdateProduct(c, db)
		})
		products.DELETE("/:id", protect, allowedTo(constants.Roles.Admin, constants.Roles.SuperAdmin), func(c *gin.Context) {
			DeleteProduct(c, db)
		})
	}
}
