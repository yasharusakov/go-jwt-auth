package config

import (
	"log"
	"os"
	"sync"
	"time"
)

type ServerConfig struct {
	AppEnv string
	Port   string
}

type PostgresConfig struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresSSLMode  string
}

type JWTConfig struct {
	JWTAccessTokenSecret  string
	JWTRefreshTokenSecret string
	JWTAccessTokenExp     time.Duration
	JWTRefreshTokenExp    time.Duration
}

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	JWT      JWTConfig
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		accessExp, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION"))
		if err != nil {
			log.Fatalf("Error parsing JWT_ACCESS_TOKEN_EXPIRATION: %v", err)
		}

		expRefresh, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION"))
		if err != nil {
			log.Fatalf("Error parsing JWT_REFRESH_TOKEN_EXPIRATION: %v", err)
		}

		cfg = &Config{
			Server: ServerConfig{
				AppEnv: os.Getenv("APP_ENV"),
				Port:   os.Getenv("API_PORT"),
			},
			Postgres: PostgresConfig{
				PostgresUser:     os.Getenv("POSTGRES_USER"),
				PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
				PostgresHost:     os.Getenv("POSTGRES_HOST"),
				PostgresPort:     os.Getenv("POSTGRES_PORT"),
				PostgresDB:       os.Getenv("POSTGRES_DB"),
				PostgresSSLMode:  os.Getenv("POSTGRES_SSLMODE"),
			},
			JWT: JWTConfig{
				JWTAccessTokenSecret:  os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
				JWTRefreshTokenSecret: os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
				JWTAccessTokenExp:     accessExp,
				JWTRefreshTokenExp:    expRefresh,
			},
		}
	})

	return cfg
}
