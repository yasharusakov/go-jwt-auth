package config

import (
	"os"
	"sync"
)

type Config struct {
	AppEnv                    string
	ApiGatewayExternalPort    string
	ClientExternalURL         string
	ApiAuthServiceInternalURL string
	ApiUserServiceInternalURL string
	JWTAccessTokenSecret      string
}

var (
	cfg  Config
	once sync.Once
)

func LoadConfigFromEnv() Config {
	return Config{
		AppEnv:                    os.Getenv("APP_ENV"),
		ApiGatewayExternalPort:    os.Getenv("API_GATEWAY_EXTERNAL_PORT"),
		ClientExternalURL:         os.Getenv("CLIENT_EXTERNAL_URL"),
		ApiAuthServiceInternalURL: os.Getenv("API_AUTH_SERVICE_INTERNAL_URL"),
		ApiUserServiceInternalURL: os.Getenv("API_USER_SERVICE_INTERNAL_URL"),
		JWTAccessTokenSecret:      os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
	}
}

func GetConfig() Config {
	once.Do(func() {
		cfg = LoadConfigFromEnv()
	})

	return cfg
}
