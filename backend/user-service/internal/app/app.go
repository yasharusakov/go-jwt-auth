package app

import (
	"context"
	"errors"
	"net/http"
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

	"google.golang.org/grpc"
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
			Msg("Error occurred while initializing postgres")
	}
	defer func() {
		logger.Log.Info().Msg("Closing postgres connection...")
		postgresGORM.Close()
	}()

	//repositories := repository.NewUserRepository(postgres)
	userRepository := repository.NewUserGormRepository(postgresGORM.DB)
	userService := service.NewUserService(userRepository)

	userHttpHandler := httpHandler.NewUserHandler(userService)
	userGrpcHandlers := grpcHandler.NewUserHandler(userService)
	routes := router.RegisterRoutes(userHttpHandler, postgresGORM.DB)

	httpServer := &server.HttpServer{}
	grpcServer := &server.GRPCServer{}

	var wg sync.WaitGroup
	serverErrors := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Info().
			Str("port", cfg.GRPCUserServiceInternalPort).
			Msg("Starting gRPC server...")

		if runErr := grpcServer.Run(cfg.GRPCUserServiceInternalPort, userGrpcHandlers); runErr != nil && !errors.Is(runErr, grpc.ErrServerStopped) {
			serverErrors <- runErr
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Info().
			Str("port", cfg.ApiUserServiceInternalPort).
			Msg("Starting HTTP server...")

		if runErr := httpServer.Run(cfg.ApiUserServiceInternalPort, routes); runErr != nil && !errors.Is(runErr, http.ErrServerClosed) {
			serverErrors <- runErr
		}
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info().Msg("Shutdown signal received")
	case err = <-serverErrors:
		logger.Log.Error().Err(err).Msg("Server error received")
		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Log.Info().Err(err).Msg("HTTP server shutdown error")
	} else {
		logger.Log.Info().Msg("HTTP server shutdown gracefully")
	}

	grpcServer.Shutdown(shutdownCtx)

	wg.Wait()
}
