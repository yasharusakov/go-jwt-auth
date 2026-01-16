package service

import (
	"auth-service/internal/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager interface {
	GenerateToken(userID string, ttl time.Duration, secret []byte) (string, error)
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	GenerateTokens(userID string) (string, string, error)

	ValidateToken(tokenStr string, secret []byte) (*jwt.RegisteredClaims, error)
	ValidateRefreshToken(tokenStr string) (*jwt.RegisteredClaims, error)
}

type manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenManager(JWT config.JWTConfig) TokenManager {
	return &manager{
		accessSecret:  []byte(JWT.JWTAccessTokenSecret),
		refreshSecret: []byte(JWT.JWTRefreshTokenSecret),
		accessTTL:     JWT.JWTAccessTokenExp,
		refreshTTL:    JWT.JWTRefreshTokenExp,
	}
}

func (m *manager) GenerateAccessToken(userID string) (string, error) {
	return m.GenerateToken(userID, m.accessTTL, m.accessSecret)
}

func (m *manager) GenerateRefreshToken(userID string) (string, error) {
	return m.GenerateToken(userID, m.refreshTTL, m.refreshSecret)
}

func (m *manager) GenerateTokens(userID string) (string, string, error) {
	accessToken, err := m.GenerateAccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := m.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (m *manager) ValidateRefreshToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	return m.ValidateToken(tokenStr, m.refreshSecret)
}

func (m *manager) GenerateToken(userID string, ttl time.Duration, secret []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func (m *manager) ValidateToken(tokenStr string, secret []byte) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}
