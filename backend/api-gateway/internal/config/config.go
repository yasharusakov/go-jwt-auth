package config

import (
	"os"
	"sync"
)

type RedisConfig struct {
	RedisInternalURL string
}

type Config struct {
	AppEnv                    string
	ApiGatewayInternalPort    string
	ClientExternalURL         string
	ApiAuthServiceInternalURL string
	ApiUserServiceInternalURL string
	JWTAccessTokenSecret      string
	RedisConfig               RedisConfig
}

var (
	cfg  Config
	once sync.Once
)

func LoadConfigFromEnv() Config {
	return Config{
		AppEnv:                    os.Getenv("APP_ENV"),
		ApiGatewayInternalPort:    os.Getenv("API_GATEWAY_INTERNAL_PORT"),
		ClientExternalURL:         os.Getenv("CLIENT_EXTERNAL_URL"),
		ApiAuthServiceInternalURL: os.Getenv("API_AUTH_SERVICE_INTERNAL_URL"),
		ApiUserServiceInternalURL: os.Getenv("API_USER_SERVICE_INTERNAL_URL"),
		JWTAccessTokenSecret:      os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
		RedisConfig: RedisConfig{
			RedisInternalURL: os.Getenv("REDIS_INTERNAL_URL"),
		},
	}
}

func GetConfig() Config {
	once.Do(func() {
		cfg = LoadConfigFromEnv()
	})

	return cfg
}
