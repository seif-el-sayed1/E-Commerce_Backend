package user

import (
	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(rg *gin.RouterGroup, db *gorm.DB, protect gin.HandlerFunc, allowedTo func(roles ...string) gin.HandlerFunc) {
	users := rg.Group("/users")
	{
		// Private routes
		users.GET("/me", protect, allowedTo(constants.Roles.User), func(c *gin.Context) {
			GetMyProfile(c, db)
		})
		users.PATCH("/me", protect, allowedTo(constants.Roles.User), func(c *gin.Context) {
			UpdateMyProfile(c, db)
		})
	}
}
