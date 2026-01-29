package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"time"
	httpHandler "user-service/internal/handler/http"
	"user-service/internal/logger"
)

func SetupRoutes(app *fiber.App, handlers httpHandler.Handlers, db *gorm.DB) {
	app.Get("/health", func(c *fiber.Ctx) error {
		logger.Log.Info().Msg("Health check passed")
		return c.SendString("OK")
	})

	app.Get("/ready", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := db.WithContext(ctx).DB()
		if err != nil {
			logger.Log.Error().Err(err).Msg("Database is not ready")
			return c.Status(fiber.StatusServiceUnavailable).SendString("Database is not ready")
		}

		if err := sqlDB.Ping(); err != nil {
			logger.Log.Error().Err(err).Msg("Database ping failed")
			return c.Status(fiber.StatusServiceUnavailable).SendString("Database is not ready")
		}

		logger.Log.Info().Msg("Ready check passed")
		return c.SendString("READY")

	})

	api := app.Group("/api")
	user := api.Group("/user")

	user.Get("/get-all", handlers.GetAllUsers)
}
