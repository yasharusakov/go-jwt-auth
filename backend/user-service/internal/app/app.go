package app

import (
	"context"
	"log"
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
)

func Run() {
	cfg := config.GetConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ========== NATS CONNECTION ==========
	// nc, err := broker.NewNATS(&cfg.NATS, "user-service")
	// if err != nil {
	//	 log.Fatal("Error occurred while connecting to NATS: ", err)
	// }
	//
	// defer func() {
	//	 log.Println("Draining NATS connection...")
	//	 if drainErr := nc.Drain(); drainErr != nil {
	//		 log.Printf("Error draining NATS connection: %v", err)
	//	 }
	// }()
	// ========== NATS END OF CONNECTION ==========

	httpServer := &server.HttpServer{}
	grpcServer := &server.GRPCServer{}

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

	// HTTP
	routes := router.RegisterRoutes(httpHandlers)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// ========== NATS HANDLERS ==========
	// natsHandler := nats.NewNATSHandler(services.User, nc)
	// natsHandler.Register()
	// log.Println("NATS handlers registered")
	// ========== END NATS HANDLERS ==========

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("gRPC server is running on port: %s", cfg.GRPCUserServicePort)
		if err := grpcServer.Run(cfg.GRPCUserServicePort, grpcHandlers); err != nil {
			errChan <- err
		}
	}()

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("HTTP server is running on port: %s", cfg.Port)
		if err := httpServer.Run(cfg.Port, routes); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err = <-errChan:
		log.Printf("Server error: %v", err)
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server shutdown gracefully")
	}
	grpcServer.GRPCServer.GracefulStop()
	log.Println("gRPC server shutdown gracefully")

	wg.Wait()
}
