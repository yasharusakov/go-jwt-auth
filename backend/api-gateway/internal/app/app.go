package app

import (
	"api-gateway/internal/config"
	"api-gateway/internal/router"
	"api-gateway/internal/server"
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

	srv := &server.HttpServer{}

	handlers := router.RegisterRoutes()

	errChan := make(chan error, 1)

	go func() {
		log.Printf("HTTP server is running on port: %s", cfg.Port)
		errChan <- srv.Run(cfg.Port, handlers)
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err := <-errChan:
		if err != nil {
			log.Fatalf("Error occurred while starting server: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Println("Server shutdown gracefully")
	}
}
