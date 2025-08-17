package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type POSTGRESQL_CONFIG struct {
	POSTGRESQL_URI string
}

type JWT_CONFIG struct {
	JWT_ACCESS_TOKEN_SECRET      string
	JWT_REFRESH_TOKEN_SECRET     string
	JWT_ACCESS_TOKEN_EXPIRATION  string
	JWT_REFRESH_TOKEN_EXPIRATION string
}

type Config struct {
	APP_ENV    string
	PORT       string
	CLIENT_URL string
	POSTGRESQL POSTGRESQL_CONFIG
	JWT        JWT_CONFIG
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func LoadConfig() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return Config{
		APP_ENV:    getEnv("APP_ENV", "development"),
		PORT:       getEnv("PORT", "8080"),
		CLIENT_URL: getEnv("CLIENT_URL", "*"),
		POSTGRESQL: POSTGRESQL_CONFIG{
			POSTGRESQL_URI: os.Getenv("POSTGRESQL_URI"),
		},
		JWT: JWT_CONFIG{
			JWT_ACCESS_TOKEN_SECRET:      os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
			JWT_REFRESH_TOKEN_SECRET:     os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
			JWT_ACCESS_TOKEN_EXPIRATION:  os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION"),
			JWT_REFRESH_TOKEN_EXPIRATION: os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION"),
		},
	}
}
