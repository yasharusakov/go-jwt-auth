package app

import (
	grpcClient "auth-service/internal/client/grpc"
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/logger"
	"auth-service/internal/repository"
	"auth-service/internal/router"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Run() {
	cfg := config.GetConfig()
	logger.Init(cfg.AppEnv)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	postgresGORM, err := storage.NewPostgresGORM(ctx, cfg.Postgres)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to Postgres.")
	}
	defer postgresGORM.Close()

	grpcUserClient, err := grpcClient.NewGRPCUserClient(cfg.GRPCUserServiceInternalURL)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to gRPC user client.")
	}
	defer grpcUserClient.Close()

	tokenRepository := repository.NewTokenRepository(postgresGORM.DB, cfg)
	tokenManager := service.NewTokenManager(cfg.JWT)
	authService := service.NewAuthService(grpcUserClient, tokenRepository, tokenManager, cfg)
	authHandler := handler.NewAuthHandler(authService, cfg)

	app := fiber.New(fiber.Config{
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           10 * time.Second,
		DisableStartupMessage: cfg.AppEnv == "production",
	})

	router.SetupRoutes(app, authHandler, postgresGORM.DB, grpcUserClient)

	serverError := make(chan error, 1)

	go func() {
		logger.Log.Info().
			Str("port", cfg.ApiAuthServiceInternalPort).
			Msg("Starting Auth Service HTTP server...")

		serverError <- app.Listen(":" + cfg.ApiAuthServiceInternalPort)
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info().Msg("Shutdown signal received.")
	case err := <-serverError:
		logger.Log.Error().Err(err).Msg("Auth Service HTTP server error.")
	}

	logger.Log.Info().Msg("Gracefully shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("Auth Service HTTP shutdown error.")
	} else {
		logger.Log.Info().Msg("Auth Service HTTP shutdown gracefully.")
	}

	logger.Log.Info().Msg("Running cleanup tasks...")
}
