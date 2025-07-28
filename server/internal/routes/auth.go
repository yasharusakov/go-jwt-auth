package routes

import (
	"net/http"
	"server/internal/handlers/auth"
	"server/internal/middlewares"
)

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/login", middlewares.CORSMiddleware(auth.Login))
	mux.HandleFunc("/api/register", middlewares.CORSMiddleware(auth.Register))
	mux.HandleFunc("/api/logout", middlewares.CORSMiddleware(auth.Logout))
	mux.HandleFunc("/api/refresh", middlewares.CORSMiddleware(auth.Refresh))
}
