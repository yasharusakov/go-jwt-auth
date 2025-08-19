package app

import (
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/database/postgresql"
	"server/internal/routes"
)

func Run() {
	cfg := config.LoadConfig()

	if err := postgresql.Connect(cfg.POSTGRESQL); err != nil {
		log.Fatal("Error occurred while connecting to the database: ", err)
	}
	defer func() {
		log.Println("Closing database connection...")
		postgresql.Pool.Close()
	}()

	router := routes.SetupRoutes()

	log.Println("Server is running on port: " + cfg.ApiPort)

	log.Fatal(http.ListenAndServe(":"+cfg.ApiPort, router))
}
