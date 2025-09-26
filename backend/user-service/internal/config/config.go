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

type NATSConfig struct {
	NatsUser     string
	NatsPassword string
	NatsHost     string
	NatsPort     string
}

type Config struct {
	Port                string
	GRPCUserServicePort string
	Postgres            PostgresConfig
	NATS                NATSConfig
}

var (
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{
			Port:                os.Getenv("API_USER_SERVICE_PORT"),
			GRPCUserServicePort: os.Getenv("GRPC_USER_SERVICE_PORT"),
			NATS: NATSConfig{
				NatsUser:     os.Getenv("NATS_USER"),
				NatsPassword: os.Getenv("NATS_PASSWORD"),
				NatsHost:     os.Getenv("NATS_HOST"),
				NatsPort:     os.Getenv("NATS_PORT"),
			},
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
