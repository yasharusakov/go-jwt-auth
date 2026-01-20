package app

import (
	"api-gateway/internal/config"
	"api-gateway/internal/router"
	"api-gateway/internal/server"
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

	handlers := router.RegisterRoutes()

	srv := &server.HttpServer{}
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("HTTP server is running on port: %s", cfg.ApiGatewayExternalPort)
		serverErrors <- srv.Run(cfg.ApiGatewayExternalPort, handlers)
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}

}
