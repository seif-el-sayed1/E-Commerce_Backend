package user

import (
	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, protect gin.HandlerFunc, allowedTo func(roles ...string) gin.HandlerFunc) {
	auth := rg.Group("/users/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			UserRegister(c, db)
		})
		auth.POST("/login", func(c *gin.Context) {
			UserLogin(c, db)
		})
		auth.POST("/verify-account", func(c *gin.Context) {
			UserVerifyAccount(c, db)
		})
		auth.POST("/verify-otp", func(c *gin.Context) {
			VerifyOTP(c, db)
		})
		auth.POST("/send-otp", func(c *gin.Context) {
			SendOTP(c, db)
		})
		auth.PATCH("/update-password", protect, allowedTo(constants.Roles.User), func(c *gin.Context) {
			UpdateLoggedUserPassword(c, db)
		})
		auth.POST("/log-out", protect, allowedTo(constants.Roles.User), func(c *gin.Context) {
			UserLogout(c, db)
		})
	}
}
