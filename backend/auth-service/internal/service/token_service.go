package service

import (
	"auth-service/internal/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	GenerateTokens(userID string) (string, string, error)
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
	return m.generateToken(userID, m.accessTTL, m.accessSecret)
}

func (m *manager) GenerateRefreshToken(userID string) (string, error) {
	return m.generateToken(userID, m.refreshTTL, m.refreshSecret)
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

func (m *manager) generateToken(userID string, ttl time.Duration, secret []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func (m *manager) ValidateRefreshToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	return m.validateToken(tokenStr, m.refreshSecret)
}

func (m *manager) validateToken(tokenStr string, secret []byte) (*jwt.RegisteredClaims, error) {
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
