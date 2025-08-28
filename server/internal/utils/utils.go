package utils

import (
	"fmt"
	"net/http"
	"os"
	"server/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int, ttl time.Duration, secret []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func GenerateTokens(userID int) (string, string, error) {
	cfg := config.LoadConfig().JWT
	accessTokenExpiration, err := time.ParseDuration(cfg.JwtAccessTokenExpiration)
	refreshTokenExpiration, err := time.ParseDuration(cfg.JwtRefreshTokenExpiration)

	accessToken, err := GenerateToken(userID, accessTokenExpiration, []byte(cfg.JwtAccessTokenSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateToken(userID, refreshTokenExpiration, []byte(cfg.JwtRefreshTokenSecret))
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

func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string, expRefresh time.Duration) {
	secure := os.Getenv("APP_ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(expRefresh),
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func RemoveRefreshTokenCookie(w http.ResponseWriter) {
	secure := os.Getenv("APP_ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}
