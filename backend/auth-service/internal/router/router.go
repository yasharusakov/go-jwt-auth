package router

import (
	"auth-service/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(handlers *handler.Handlers) http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/api/auth/register", handlers.Auth.Register).Methods("POST")
	m.HandleFunc("/api/auth/login", handlers.Auth.Login).Methods("POST")
	m.HandleFunc("/api/auth/refresh", handlers.Auth.Refresh).Methods("GET")
	m.HandleFunc("/api/auth/logout", handlers.Auth.Logout).Methods("POST")

	return m
}
