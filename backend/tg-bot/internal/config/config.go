package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Token                       string
	AdminID                     int
	GrpcUserServiceInternalAddr string
}

var (
	cfg  Config
	once sync.Once
)

func loadConfigFromEnv() Config {

	adminID, err := strconv.Atoi(os.Getenv("TG_BOT_ADMIN_ID"))
	if err != nil {
		panic("Invalid TG_BOT_ADMIN_ID: " + err.Error())
	}

	return Config{
		Token:                       os.Getenv("TG_BOT_TOKEN"),
		AdminID:                     adminID,
		GrpcUserServiceInternalAddr: os.Getenv("GRPC_USER_SERVICE_INTERNAL_URL"),
	}
}

func GetConfig() Config {
	once.Do(func() {
		cfg = loadConfigFromEnv()
	})
	return cfg
}
