package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/admin"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/category"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/product"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/subcategory"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/user"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/variant"
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

	MigrateModels()
}

func MigrateModels() {
	err := DB.AutoMigrate(
		admin.Admin{},
		user.User{},
		category.Category{},
		subcategory.SubCategory{},
		variant.Variant{},
		product.Product{},
	)
	if err != nil {
		log.Fatal(constants.Error("Failed to migrate database: "), err)
	}

	log.Println(constants.Success(" 🚀 Database migrated successfully"))
}
