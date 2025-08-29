package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/internal/database"
	"server/internal/handler"
	"server/internal/repository"
	"server/internal/routes"
	"server/internal/service"
	"syscall"
	"time"
)

func Run() {
	ctx := context.Background()

	postgres, err := database.NewPostgres(ctx, database.PostgresConfig{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresSSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
	})
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

	router := routes.SetupRoutes(handlers)

	server := &Server{}
	apiPort := os.Getenv("API_PORT")
	go func() {
		if err = server.Run(apiPort, router); err != nil {
			log.Fatal("Error occurred while starting server: ", err.Error())
		}
	}()

	log.Println("Server is running on port: " + apiPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Server is shutting down...")

	if err = server.Shutdown(ctx); err != nil {
		log.Fatal("Error occurred on server shutting down: ", err.Error())
	}

	log.Println("Server has been shut down.")
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
