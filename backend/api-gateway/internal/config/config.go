package config

import (
	"os"
	"sync"
)

type Config struct {
	Port                 string
	ClientURL            string
	JWTAccessTokenSecret string
	ApiAuthServiceURL    string
	ApiUserServiceURL    string
}

var (
	cfg  Config
	once sync.Once
)

func LoadConfigFromEnv() Config {
	return Config{
		Port:                 os.Getenv("API_GATEWAY_PORT"),
		ClientURL:            os.Getenv("CLIENT_URL"),
		JWTAccessTokenSecret: os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
		ApiAuthServiceURL:    os.Getenv("API_AUTH_SERVICE_URL"),
		ApiUserServiceURL:    os.Getenv("API_USER_SERVICE_URL"),
	}
}

func GetConfig() Config {
	once.Do(func() {
		cfg = LoadConfigFromEnv()
	})

	return cfg
}
