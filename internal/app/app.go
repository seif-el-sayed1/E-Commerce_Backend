package app

import (
	"seif-el-sayed1/E-Commerce_Backend.git/internal/admin"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/database"
)

/**
 this Package initialization entry point.
 This file aggregates initialization functions from different packages
 across the application and executes them when the server starts.
**/

func Run() {
	config.LoadEnv()
	database.ConnectToDb()
	admin.SeedSuperAdmin(database.DB)
}
