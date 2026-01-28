package app

import (
	"api-gateway/internal/cache"
	"api-gateway/internal/config"
	"api-gateway/internal/logger"
	"api-gateway/internal/router"
	"context"
	"github.com/gofiber/fiber/v2"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg := config.GetConfig()
	logger.Init(cfg.AppEnv)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	redisCache := cache.NewRedisCache(cfg.RedisConfig)
	defer redisCache.Close()

	app := fiber.New(fiber.Config{
		WriteTimeout:          10 * time.Second,
		ReadTimeout:           10 * time.Second,
		IdleTimeout:           10 * time.Second,
		DisableStartupMessage: cfg.AppEnv == "production",
	})

	router.SetupRoutes(app, redisCache, cfg)

	go func() {
		logger.Log.Info().
			Str("port", ":"+cfg.ApiGatewayInternalPort).
			Msg("Starting API Gateway server...")

		if err := app.Listen(":" + cfg.ApiGatewayInternalPort); err != nil {
			logger.Log.Panic().Err(err)
		}
	}()

	<-ctx.Done()
	logger.Log.Info().Msg("Gracefully shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("API Gateway shutdown error.")
	} else {
		logger.Log.Info().Msg("API Gateway stopped gracefully.")
	}

	logger.Log.Info().Msg("Running cleanup tasks...")

	logger.Log.Info().Msg("API Gateway was successful shutdown.")
}
