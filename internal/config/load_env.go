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
	BASE_URL            string
	JWTSecret           string
	JWTExpiration       string
	SuperAdminEmail     string
	SuperAdminPassword  string
	SuperAdminFirstName string
	SuperAdminLastName  string
	EmailPass           string
	EmailUser           string
	AdminUrl            string
	AdminLogin          string
	AdminVerify         string
	AdminForgotPass     string
	AppName             string
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
		BASE_URL:            os.Getenv("BASE_URL"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTExpiration:       os.Getenv("JWT_EXPIRATION"),
		SuperAdminEmail:     os.Getenv("SUPER_ADMIN_EMAIL"),
		SuperAdminPassword:  os.Getenv("SUPER_ADMIN_PASSWORD"),
		SuperAdminFirstName: os.Getenv("SUPER_ADMIN_FIRST_NAME"),
		SuperAdminLastName:  os.Getenv("SUPER_ADMIN_LAST_NAME"),
		EmailPass:           os.Getenv("EMAIL_PASS"),
		EmailUser:           os.Getenv("EMAIL_USER"),
		AdminUrl:            os.Getenv("ADMIN_URL"),
		AdminVerify:         os.Getenv("ADMIN_VERIFY"),
		AdminLogin:          os.Getenv("ADMIN_LOGIN"),
		AdminForgotPass:     os.Getenv("ADMIN_FORGOT_PASS"),
		AppName:             os.Getenv("APP_NAME"),
	}

	log.Println(constants.Success("✅ Environment variables loaded"))
}
