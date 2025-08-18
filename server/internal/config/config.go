package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
	APP_ENV     string
	API_PORT    string
	CLIENT_PORT string
	CLIENT_URL  string
	POSTGRESQL  POSTGRESQL_CONFIG
	JWT         JWT_CONFIG
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

// FIXME: The environment file is loaded every time this function is called and should be loaded only once
func LoadConfig() Config {
	// FIXME: It prints an error if the .env.development file is not found, but I have others like prod or docker.
	if err := godotenv.Load(".env.development"); err != nil {
		log.Printf("Error loading .env.development file: %v", err)
	}

	return Config{
		APP_ENV:     getEnv("APP_ENV", "development"),
		API_PORT:    getEnv("API_PORT", "8080"),
		CLIENT_PORT: getEnv("CLIENT_PORT", "3000"),
		CLIENT_URL:  getEnv("CLIENT_URL", "*"),
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
