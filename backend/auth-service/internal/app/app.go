package app

import (
	grpcClient "auth-service/internal/client/grpc"
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/logger"
	"auth-service/internal/repository"
	"auth-service/internal/router"
	"auth-service/internal/server"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg := config.GetConfig()
	logger.Init(cfg.AppEnv)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//postgres, err := storage.NewPostgres(ctx, cfg.Postgres)
	postgresGORM, err := storage.NewPostgresGORM(ctx, cfg.Postgres)
	if err != nil {
		logger.Log.Fatal().
			Err(err).
			Msg("Failed to connect to Postgres")
	}
	defer func() {
		logger.Log.Info().Msg("Closing postgres connection...")
		postgresGORM.Close()
	}()

	grpcUserClient, err := grpcClient.NewGRPCUserClient(cfg.GRPCUserServiceInternalURL)
	if err != nil {
		logger.Log.Fatal().
			Err(err).
			Msg("Failed to connect to gRPC user client")
	}
	defer func() {
		logger.Log.Info().Msg("Closing gRPC user client...")
		grpcUserClient.Close()
	}()

	tokenRepo := repository.NewTokenRepository(postgresGORM.DB, cfg)
	tokenManager := service.NewTokenManager(cfg.JWT)
	authService := service.NewAuthService(grpcUserClient, tokenRepo, tokenManager, cfg)
	authHandler := handler.NewAuthHandler(authService, cfg)
	routes := router.RegisterRoutes(authHandler, postgresGORM.DB, grpcUserClient)

	httpServer := &server.HttpServer{}
	serverErrors := make(chan error, 1)

	go func() {
		logger.Log.Info().
			Str("port", cfg.ApiAuthServiceInternalPort).
			Msg("Starting HTTP server...")

		serverErrors <- httpServer.Run(cfg.ApiAuthServiceInternalPort, routes)
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info().Msg("Shutdown signal received")
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal().
				Err(err).
				Msg("Server error")
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Server shutdown error")
	} else {
		logger.Log.Info().Msg("Server shutdown gracefully")
	}
}
