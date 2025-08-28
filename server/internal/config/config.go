package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresSSLMode  string
}

type JwtConfig struct {
	JwtAccessTokenSecret      string
	JwtRefreshTokenSecret     string
	JwtAccessTokenExpiration  string
	JwtRefreshTokenExpiration string
}

type Config struct {
	AppEnv     string
	ApiPort    string
	ClientPort string
	ClientUrl  string
	Postgres   PostgresConfig
	JWT        JwtConfig
}

var (
	config *Config
	once   sync.Once
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func loadEnv() {
	if os.Getenv("APP_ENV") == "docker" {
		log.Println("Docker environment detected, loaded successfully.")
	} else {
		if err := godotenv.Load(".env.development"); err != nil {
			log.Fatalf("Error loading .env.development file: %v", err)
		} else {
			log.Println("Development environment detected, loaded successfully.")
			return
		}
	}
}

func LoadConfig() *Config {
	once.Do(func() {
		loadEnv()

		config = &Config{
			AppEnv:     getEnv("APP_ENV", "development"),
			ApiPort:    getEnv("API_PORT", "8080"),
			ClientPort: getEnv("CLIENT_PORT", "3000"),
			ClientUrl:  getEnv("CLIENT_URL", "*"),
			Postgres: PostgresConfig{
				PostgresUser:     os.Getenv("POSTGRES_USER"),
				PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
				PostgresHost:     os.Getenv("POSTGRES_HOST"),
				PostgresPort:     os.Getenv("POSTGRES_PORT"),
				PostgresDB:       os.Getenv("POSTGRES_DB"),
				PostgresSSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
			},
			JWT: JwtConfig{
				JwtAccessTokenSecret:      os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
				JwtRefreshTokenSecret:     os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
				JwtAccessTokenExpiration:  os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION"),
				JwtRefreshTokenExpiration: os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION"),
			},
		}
	})

	return config
}
