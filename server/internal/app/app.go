package app

import (
	"context"
	"log"
	"os/signal"
	"server/internal/config"
	"server/internal/handler"
	"server/internal/repository"
	"server/internal/router"
	"server/internal/server"
	"server/internal/service"
	"server/internal/storage"
	"syscall"
	"time"
)

func Run() {
	cfg := config.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	postgres, err := storage.NewPostgres(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal("Error occurred while initializing postgres: ", err)
	}
	defer func() {
		log.Println("Closing initial postgres connection...")
		postgres.Close()
	}()

	repositories := repository.NewRepositories(postgres)
	services := service.NewServices(repositories)
	handlers := handler.NewHandlers(services)
	r := router.SetupRoutes(handlers)

	srv := &server.Server{}
	go func() {
		if err = srv.Run(cfg.Server.Port, r); err != nil {
			log.Fatal("Error occurred while starting server: ", err.Error())
		}
	}()

	log.Println("Server is running on port: " + cfg.Server.Port)

	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Println("Server shutdown gracefully")
	}
}
