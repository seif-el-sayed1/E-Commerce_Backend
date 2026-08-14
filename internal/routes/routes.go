package routes

import (
	"seif-el-sayed1/E-Commerce_Backend.git/internal/admin"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AppRoutes(server *gin.Engine, db *gorm.DB) {
	api := server.Group(config.Env.BASE_URL)

	admin.AdminAuthRoutes(api, db, middlewares.Protect(db))
	admin.AdminRoutes(api, db, middlewares.Protect(db), middlewares.AllowedTo)
}
