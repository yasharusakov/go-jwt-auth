package middleware

import (
	"api-gateway/internal/config"
	"api-gateway/internal/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Log.Warn().
				Str("path", c.Path()).
				Str("ip", c.IP()).
				Msg("access token not found")

			return c.SendStatus(fiber.StatusUnauthorized)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token
		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(cfg.JWTAccessTokenSecret), nil
		})

		if err != nil || !token.Valid {
			logger.Log.Warn().
				Err(err).
				Str("path", c.Path()).
				Str("ip", c.IP()).
				Msg("invalid or expired token")

			return c.SendStatus(fiber.StatusUnauthorized)
		}

		_, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			logger.Log.Warn().
				Str("path", c.Path()).
				Str("ip", c.IP()).
				Msg("invalid claims")

			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}
