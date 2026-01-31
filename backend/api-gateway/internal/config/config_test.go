package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setTestEnv(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("API_GATEWAY_INTERNAL_PORT", "8080")
	t.Setenv("CLIENT_EXTERNAL_URL", "https://example.com")
	t.Setenv("API_AUTH_SERVICE_INTERNAL_URL", "http://auth:3001")
	t.Setenv("API_USER_SERVICE_INTERNAL_URL", "http://user:3002")
	t.Setenv("JWT_ACCESS_TOKEN_SECRET", "secret")
	t.Setenv("REDIS_INTERNAL_URL", "redis://localhost:6379")
}

func TestLoadConfigFromEnv(t *testing.T) {
	setTestEnv(t)

	cfg := LoadConfigFromEnv()

	assert.Equal(t, "test", cfg.AppEnv)
	assert.Equal(t, "8080", cfg.ApiGatewayInternalPort)
	assert.Equal(t, "https://example.com", cfg.ClientExternalURL)
	assert.Equal(t, "http://auth:3001", cfg.ApiAuthServiceInternalURL)
	assert.Equal(t, "http://user:3002", cfg.ApiUserServiceInternalURL)
	assert.Equal(t, "secret", cfg.JWTAccessTokenSecret)
	assert.Equal(t, "redis://localhost:6379", cfg.RedisConfig.RedisInternalURL)
}

func TestLoadConfigFromEnv_Empty(t *testing.T) {
	cfg := LoadConfigFromEnv()

	assert.Empty(t, cfg.AppEnv)
	assert.Empty(t, cfg.ApiGatewayInternalPort)
}
