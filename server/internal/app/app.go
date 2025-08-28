package app

import (
	"context"
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/database"
	"server/internal/handler"
	"server/internal/repository"
	"server/internal/routes"
	"server/internal/service"
)

func Run() {
	ctx := context.Background()
	cfg := config.LoadConfig()

	postgres, err := database.NewPostgres(ctx, cfg.Postgres)
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

	log.Println("Server is running on port: " + cfg.ApiPort)

	log.Fatal(http.ListenAndServe(":"+cfg.ApiPort, router))
}
