package app

import (
	"auth-service/internal/broker"
	"auth-service/internal/client"
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/router"
	"auth-service/internal/server"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg := config.GetConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ========== NATS CONNECTION ==========
	nc, err := broker.NewNATS(&cfg.NATS, "user-service")
	if err != nil {
		log.Fatal("Error occurred while connecting to NATS: ", err)
	}

	defer func() {
		log.Println("Draining NATS connection...")
		if drainErr := nc.Drain(); drainErr != nil {
			log.Printf("Error draining NATS connection: %v", err)
		}
	}()
	// ========== NATS END OF CONNECTION ==========

	srv := &server.Server{}

	postgres, err := storage.NewPostgres(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal("Error occurred while initializing postgres: ", err)
	}
	defer func() {
		log.Println("Closing initial postgres connection...")
		postgres.Close()
	}()

	repositories := repository.NewRepositories(postgres)
	userClient := client.NewUserClient(nc)
	services := service.NewServices(userClient, repositories)
	handlers := handler.NewHandlers(services)
	routes := router.RegisterRoutes(handlers)

	errChan := make(chan error, 1)

	go func() {
		errChan <- srv.Run(cfg.Port, routes)
	}()

	log.Printf("Server is running on port: %s", cfg.Port)

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err = <-errChan:
		if err != nil {
			log.Fatalf("Error occurred while starting server: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Println("Server shutdown gracefully")
	}
}
