package app

import (
	"seif-el-sayed1/E-Commerce_Backend.git/internal/admin"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"
)

/**
 Package initialization entry point.
 This file aggregates initialization functions from different packages
 across the application and executes them when the server starts.
**/

func Run() {
	config.LoadEnv()
	config.ConnectToDb()
	admin.SeedSuperAdmin(config.DB, config.Env.SuperAdminEmail, config.Env.SuperAdminPassword)
}
