package admin

import (
	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"

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
			AdminChangePassword(c, db)
		})
		auth.POST("/verify-account", protect, allowedTo(constants.Roles.SuperAdmin, constants.Roles.Admin), func(c *gin.Context) {
			VerifyAdminAccount(c, db)
		})
		auth.POST("forget-password", func(c *gin.Context) {
			AdminForgetPassword(c, db)
		})
		auth.PATCH("/reset-password/:token", func(c *gin.Context) {
			AdminResetPassword(c, db)
		})
	}

}
