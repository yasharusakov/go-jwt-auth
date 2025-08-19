package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type PostgresqlConfig struct {
	PostgresqlUri string
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
	POSTGRESQL PostgresqlConfig
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
			POSTGRESQL: PostgresqlConfig{
				PostgresqlUri: os.Getenv("POSTGRESQL_URI"),
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
