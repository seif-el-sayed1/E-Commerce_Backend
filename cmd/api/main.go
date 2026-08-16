package main

import (
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/app"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()
	app.Run(server)

	port := config.Env.Port
	if port == "" {
		port = "8080" // Default port if not specified in .env
	}

	server.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Hello World",
		})
	})

	server.NoRoute(func(c *gin.Context) {
		apiErr := utils.NewApiError(
			"Route not found: "+c.Request.Method+" "+c.Request.URL.Path,
			http.StatusNotFound,
		)
		c.JSON(apiErr.StatusCode, gin.H{
			"success": apiErr.Success,
			"status":  apiErr.Status,
			"message": apiErr.Message,
		})
	})

	server.Run(":" + port)
}
