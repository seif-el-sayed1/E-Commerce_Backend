package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"
)

type Config struct {
	DatabaseURL         string
	Port                string
	JWTSecret           string
	JWTExpiration       string
	SuperAdminEmail     string
	SuperAdminPassword  string
	SuperAdminFirstName string
	SuperAdminLastName  string
}

var Env Config

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(constants.Warning("⚠️ No .env file found, reading from system environment variables"))
	}

	requiredVars := []string{
		"DATABASE_URL",
		"PORT",
		"JWT_SECRET",
		"JWT_EXPIRATION",
	}

	for _, key := range requiredVars {
		if os.Getenv(key) == "" {
			log.Fatal(constants.Error("❌ Missing required environment variable: "), key)
		}
	}

	Env = Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		Port:                os.Getenv("PORT"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTExpiration:       os.Getenv("JWT_EXPIRATION"),
		SuperAdminEmail:     os.Getenv("SUPER_ADMIN_EMAIL"),
		SuperAdminPassword:  os.Getenv("SUPER_ADMIN_PASSWORD"),
		SuperAdminFirstName: os.Getenv("SUPER_ADMIN_FIRST_NAME"),
		SuperAdminLastName:  os.Getenv("SUPER_ADMIN_LAST_NAME"),
	}

	log.Println(constants.Success("✅ Environment variables loaded"))
}
