package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"user-service/internal/config"
	grpcHandler "user-service/internal/handler/grpc"
	httpHandler "user-service/internal/handler/http"
	"user-service/internal/repository"
	"user-service/internal/router"
	"user-service/internal/server"
	"user-service/internal/service"
	"user-service/internal/storage"

	"google.golang.org/grpc"
)

func Run() {
	cfg := config.GetConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	postgres, err := storage.NewPostgres(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal("Error occurred while initializing postgres: ", err)
	}
	defer func() {
		log.Println("Closing postgres connection...")
		postgres.Close()
	}()

	repositories := repository.NewUserRepository(postgres)
	services := service.NewUserService(repositories)

	httpHandlers := httpHandler.NewUserHandler(services)
	grpcHandlers := grpcHandler.NewUserHandler(services)
	routes := router.RegisterRoutes(httpHandlers)

	httpServer := &server.HttpServer{}
	grpcServer := &server.GRPCServer{}

	var wg sync.WaitGroup
	serverErrors := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("gRPC server is running on port: %s", cfg.GRPCUserServiceInternalPort)

		if runErr := grpcServer.Run(cfg.GRPCUserServiceInternalPort, grpcHandlers); runErr != nil && !errors.Is(runErr, grpc.ErrServerStopped) {
			serverErrors <- runErr
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("HTTP server is running on port: %s", cfg.ApiUserServiceInternalPort)

		if runErr := httpServer.Run(cfg.ApiUserServiceInternalPort, routes); runErr != nil && !errors.Is(runErr, http.ErrServerClosed) {
			serverErrors <- runErr
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err = <-serverErrors:
		log.Printf("Server error: %v", err)
		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server shutdown gracefully")
	}

	grpcServer.Shutdown(shutdownCtx)

	wg.Wait()
}
