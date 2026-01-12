package util

import (
	"auth-service/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, ttl time.Duration, secret []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func GenerateTokens(userID string, jwtCfg config.JWTConfig) (string, string, error) {
	accessToken, err := GenerateToken(userID, jwtCfg.JWTAccessTokenExp, []byte(jwtCfg.JWTAccessTokenSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateToken(userID, jwtCfg.JWTRefreshTokenExp, []byte(jwtCfg.JWTRefreshTokenSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenStr string, secret []byte) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
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

func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string, exp time.Duration, isProd bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(exp),
		HttpOnly: true,
		Path:     "/",
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
	})
}

func RemoveRefreshTokenCookie(w http.ResponseWriter, isProd bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
	})
}
