package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"
)

var DB *gorm.DB

func ConnectToDb() {
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(constants.Error("❌ Failed to connect to database: "), err)
	}

	log.Println(constants.Success("🚀 Database connected successfully"))
	DB = db

}
