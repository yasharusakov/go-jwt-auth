package middleware

import (
	"api-gateway/internal/cache"

	"github.com/gofiber/fiber/v2"
)

func RateLimit(redisCache cache.RedisCache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		allowed, err := redisCache.CheckRateLimit(c.Context(), clientIP)
		if err != nil {
			// Redis error - pass
			return c.Next()
		}

		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate limit exceeded",
			})
		}

		return c.Next()
	}
}
