package app

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"user-service/internal/config"
	grpcHandler "user-service/internal/handler/grpc"
	"user-service/internal/repository"
	"user-service/internal/server"
	"user-service/internal/service"
	"user-service/internal/storage"
)

func Run() {
	cfg := config.GetConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ========== NATS CONNECTION ==========
	//nc, err := broker.NewNATS(&cfg.NATS, "user-service")
	//if err != nil {
	//	log.Fatal("Error occurred while connecting to NATS: ", err)
	//}
	//
	//defer func() {
	//	log.Println("Draining NATS connection...")
	//	if drainErr := nc.Drain(); drainErr != nil {
	//		log.Printf("Error draining NATS connection: %v", err)
	//	}
	//}()
	// ========== NATS END OF CONNECTION ==========

	//httpServer := &server.HttpServer{}
	grpcServer := &server.GRPCServer{}

	postgres, err := storage.NewPostgres(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal("Error occurred while initializing postgres: ", err)
	}
	defer func() {
		log.Println("Closing initial postgres connection...")
		postgres.Close()
	}()

	//repositories := repository.NewUserRepository(postgres)
	//services := service.NewUserService(repositories)
	//handlers := httpHandler.NewUserHandler(services)
	//routes := router.RegisterRoutes(handlers)

	repositories := repository.NewUserRepository(postgres)
	services := service.NewUserService(repositories)
	handlers := grpcHandler.NewUserHandler(services)

	// ========== NATS HANDLERS ==========
	//natsHandler := nats.NewNATSHandler(services.User, nc)
	//natsHandler.Register()
	//log.Println("NATS handlers registered")
	// ========== END NATS HANDLERS ==========

	errChan := make(chan error, 1)

	go func() {
		//errChan <- httpServer.Run(cfg.Port, routes)
		errChan <- grpcServer.Run(cfg.GRPCUserServicePort, handlers)
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

	//shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	//if err = httpServer.Shutdown(shutdownCtx); err != nil {
	//	log.Printf("Server shutdown error: %v", err)
	//} else {
	//	log.Println("Server shutdown gracefully")
	//}
	grpcServer.GRPCServer.GracefulStop()
}
