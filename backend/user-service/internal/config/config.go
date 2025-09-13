package config

import (
	"os"
	"sync"
)

type PostgresConfig struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresSSLMode  string
}

type Config struct {
	Port     string
	Postgres PostgresConfig
}

var (
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{
			Port: os.Getenv("API_USER_SERVICE_PORT"),
			Postgres: PostgresConfig{
				PostgresUser:     os.Getenv("DB_USER_POSTGRES_USER"),
				PostgresPassword: os.Getenv("DB_USER_POSTGRES_PASSWORD"),
				PostgresHost:     os.Getenv("DB_USER_POSTGRES_HOST"),
				PostgresPort:     os.Getenv("DB_USER_POSTGRES_INTERNAL_PORT"),
				PostgresDB:       os.Getenv("DB_USER_POSTGRES_DB"),
				PostgresSSLMode:  os.Getenv("DB_USER_POSTGRES_SSLMODE"),
			},
		}
	})

	return cfg
}
