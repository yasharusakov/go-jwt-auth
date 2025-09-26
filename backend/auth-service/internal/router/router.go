package router

import (
	"auth-service/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(handlers handler.AuthHandler) http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/api/auth/register", handlers.Register).Methods("POST")
	m.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")
	m.HandleFunc("/api/auth/refresh", handlers.Refresh).Methods("GET")
	m.HandleFunc("/api/auth/logout", handlers.Logout).Methods("POST")

	return m
}
