package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

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
	AppEnv              string
	Port                string
	ApiUserServiceURL   string
	GRPCUserServicePort string
	GRPCUserServiceURL  string
	Postgres            PostgresConfig
	JWT                 JWTConfig
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfigFromEnv() (*Config, error) {
	accessExp, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT_ACCESS_TOKEN_EXPIRATION: %w", err)
	}

	expRefresh, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT_REFRESH_TOKEN_EXPIRATION: %w", err)
	}

	return &Config{
		AppEnv:              os.Getenv("APP_ENV"),
		Port:                os.Getenv("API_AUTH_SERVICE_PORT"),
		ApiUserServiceURL:   os.Getenv("API_USER_SERVICE_URL"),
		GRPCUserServicePort: os.Getenv("GRPC_USER_SERVICE_PORT"),
		GRPCUserServiceURL:  os.Getenv("GRPC_USER_SERVICE_URL"),
		Postgres: PostgresConfig{
			PostgresUser:     os.Getenv("DB_AUTH_POSTGRES_USER"),
			PostgresPassword: os.Getenv("DB_AUTH_POSTGRES_PASSWORD"),
			PostgresHost:     os.Getenv("DB_AUTH_POSTGRES_HOST"),
			PostgresPort:     os.Getenv("DB_AUTH_POSTGRES_INTERNAL_PORT"),
			PostgresDB:       os.Getenv("DB_AUTH_POSTGRES_DB"),
			PostgresSSLMode:  os.Getenv("DB_AUTH_POSTGRES_SSLMODE"),
		},
		JWT: JWTConfig{
			JWTAccessTokenSecret:  os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
			JWTRefreshTokenSecret: os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
			JWTAccessTokenExp:     accessExp,
			JWTRefreshTokenExp:    expRefresh,
		},
	}, err
}

func GetConfig() *Config {
	once.Do(func() {
		var err error
		cfg, err = LoadConfigFromEnv()
		if err != nil {
			log.Fatal(err)
		}
	})

	return cfg
}
