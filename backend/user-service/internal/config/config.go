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
	ApiUserServiceInternalPort  string
	GRPCUserServiceInternalPort string
	Postgres                    PostgresConfig
}

var (
	cfg  Config
	once sync.Once
)

func LoadConfigFromEnv() Config {
	return Config{
		ApiUserServiceInternalPort:  os.Getenv("API_USER_SERVICE_INTERNAL_PORT"),
		GRPCUserServiceInternalPort: os.Getenv("GRPC_USER_SERVICE_INTERNAL_PORT"),
		Postgres: PostgresConfig{
			PostgresUser:     os.Getenv("DB_USER_POSTGRES_USER"),
			PostgresPassword: os.Getenv("DB_USER_POSTGRES_PASSWORD"),
			PostgresHost:     os.Getenv("DB_USER_POSTGRES_HOST"),
			PostgresPort:     os.Getenv("DB_USER_POSTGRES_INTERNAL_PORT"),
			PostgresDB:       os.Getenv("DB_USER_POSTGRES_DB"),
			PostgresSSLMode:  os.Getenv("DB_USER_POSTGRES_SSLMODE"),
		},
	}
}

func GetConfig() Config {
	once.Do(func() {
		cfg = LoadConfigFromEnv()
	})

	return cfg
}
