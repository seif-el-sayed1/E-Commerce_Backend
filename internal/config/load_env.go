package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"
)

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

	log.Println(constants.Success("✅ Environment variables loaded"))
}
