package router

import (
	grpcClient "auth-service/internal/client/grpc"
	"auth-service/internal/handler"
	"auth-service/internal/logger"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(
	app *fiber.App,
	handlers handler.AuthHandler,
	db *gorm.DB,
	grpcUserClient grpcClient.UserService,
) {
	app.Get("/health", func(c *fiber.Ctx) error {
		logger.Log.Info().Msg("Health check passed")
		return c.SendStatus(http.StatusOK)
	})

	app.Get("/ready", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := db.WithContext(ctx).DB()
		if err != nil {
			logger.Log.Error().Err(err).Msg("Database is not ready")
			return c.Status(http.StatusServiceUnavailable).SendString("Database is not ready")
		}

		if err := sqlDB.Ping(); err != nil {
			logger.Log.Error().Err(err).Msg("Database ping failed")
			return c.Status(http.StatusServiceUnavailable).SendString("Database is not ready")
		}

		if err := grpcUserClient.Ping(ctx); err != nil {
			logger.Log.Error().Err(err).Msg("gRPC user client ping failed")
			return c.Status(http.StatusServiceUnavailable).SendString("gRPC user client is not ready")
		}

		logger.Log.Info().Msg("Ready check passed")

		return c.SendStatus(http.StatusOK)
	})

	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	api := app.Group("/api")
	auth := api.Group("/auth")

	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
	auth.Get("/refresh", handlers.Refresh)
	auth.Post("/logout", handlers.Logout)
}
