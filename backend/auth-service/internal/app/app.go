package app

import (
	grpcClient "auth-service/internal/client/grpc"
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/router"
	"auth-service/internal/server"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg := config.GetConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	postgres, err := storage.NewPostgres(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to init postgres: ", err)
	}
	defer func() {
		log.Println("Closing postgres connection...")
		postgres.Close()
	}()

	grpcUserClient, err := grpcClient.NewGRPCUserClient(cfg.GRPCUserServiceInternalURL)
	if err != nil {
		log.Fatal("Failed to create gRPC user client: ", err)
	}
	defer func() {
		log.Println("closing gRPC user client...")
		grpcUserClient.Close()
	}()

	tokenRepo := repository.NewTokenRepository(postgres, cfg)
	tokenManager := service.NewTokenManager(cfg.JWT)
	authService := service.NewAuthService(grpcUserClient, tokenRepo, tokenManager, cfg)
	authHandler := handler.NewAuthHandler(authService, cfg)
	routes := router.RegisterRoutes(authHandler, postgres, grpcUserClient)

	httpServer := &server.HttpServer{}
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("HTTP server is running on port: %s", cfg.ApiAuthServiceInternalPort)
		serverErrors <- httpServer.Run(cfg.ApiAuthServiceInternalPort, routes)
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
