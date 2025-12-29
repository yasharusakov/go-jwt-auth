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

	srv := &server.HttpServer{}

	handlers := router.RegisterRoutes()

	errChan := make(chan error, 1)

	go func() {
		log.Printf("HTTP server is running on port: %s", cfg.Port)
		errChan <- srv.Run(cfg.Port, handlers)
	}()

	handleRunErr := func(err error) {
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		} else {
			log.Println("Stopped cleanly")
		}
	}

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err := <-errChan:
		handleRunErr(err)
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Println("Server shutdown gracefully")
	}

	handleRunErr(<-errChan)
}
