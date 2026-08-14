package main

import (
	"net/http"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/app"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"

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

	server.Run(":" + port)
}
