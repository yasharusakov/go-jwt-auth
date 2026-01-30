package app

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"user-service/internal/config"
	grpcHandler "user-service/internal/handler/grpc"
	httpHandler "user-service/internal/handler/http"
	"user-service/internal/logger"
	"user-service/internal/repository"
	"user-service/internal/router"
	"user-service/internal/server"
	"user-service/internal/service"
	"user-service/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func Run() {
	cfg := config.GetConfig()
	logger.Init(cfg.AppEnv)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	postgresGORM, err := storage.NewPostgresGORM(ctx, cfg.Postgres)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error occurred while initializing postgres")
	}
	defer postgresGORM.Close()

	userRepository := repository.NewUserGormRepository(postgresGORM.DB)
	userService := service.NewUserService(userRepository)

	userHttpHandlers := httpHandler.NewUserHandler(userService)
	userGrpcHandlers := grpcHandler.NewUserHandler(userService)

	app := fiber.New(fiber.Config{
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           10 * time.Second,
		DisableStartupMessage: cfg.AppEnv == "production",
	})

	router.SetupRoutes(app, userHttpHandlers, postgresGORM.DB)

	grpcServer := &server.GRPCServer{}

	var wg sync.WaitGroup
	serverErrors := make(chan error, 2)

	// gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Info().
			Str("port", cfg.GRPCUserServiceInternalPort).
			Msg("Starting gRPC server...")

		serverErrors <- grpcServer.Run(cfg.GRPCUserServiceInternalPort, userGrpcHandlers)
	}()

	// HTTP server (Fiber)
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Info().
			Str("port", cfg.ApiUserServiceInternalPort).
			Msg("Starting HTTP server...")

		serverErrors <- app.Listen(":" + cfg.ApiUserServiceInternalPort)
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info().Msg("Shutdown signal received.")
	case err = <-serverErrors:
		logger.Log.Error().Err(err).Msg("Server error received.")
		stop()
	}

	logger.Log.Info().Msg("Gracefully shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("User Service HTTP server shutdown error.")
	} else {
		logger.Log.Info().Msg("User Service HTTP server stopped gracefully.")
	}

	grpcServer.Shutdown(shutdownCtx)

	wg.Wait()

	logger.Log.Info().Msg("Running cleanup tasks...")
}
