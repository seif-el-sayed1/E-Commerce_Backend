package admin

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, protect gin.HandlerFunc, allowedTo func(roles ...string) gin.HandlerFunc) {
	auth := rg.Group("/admins/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			Login(c, db)
		})
		auth.PATCH("/change-password", protect, func(c *gin.Context) {
			AdminChangePassowrd(c, db)
		})
	}

}
