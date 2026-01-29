package middleware

import (
	"api-gateway/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS(cfg config.Config) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.ClientExternalURL,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
	})
}
