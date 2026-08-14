package admin

import (
	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(rg *gin.RouterGroup, db *gorm.DB, protect gin.HandlerFunc, allowedTo func(roles ...string) gin.HandlerFunc) {
	adminRoutes := rg.Group("/admins")
	{
		adminRoutes.GET("/profile", protect, allowedTo(constants.Roles.SuperAdmin), func(c *gin.Context) {
			GetMyProfile(c, db)
		})
	}

}
