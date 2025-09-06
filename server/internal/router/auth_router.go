package router

import (
	"net/http"
	"server/internal/handler"
	"server/internal/middleware"
)

func RegisterAuthRoutes(mux *http.ServeMux, authHandler handler.AuthHandler) {
	mux.HandleFunc("/api/login", middleware.CORSMiddleware(authHandler.Login))
	mux.HandleFunc("/api/register", middleware.CORSMiddleware(authHandler.Register))
	mux.HandleFunc("/api/logout", middleware.CORSMiddleware(authHandler.Logout))
	mux.HandleFunc("/api/refresh", middleware.CORSMiddleware(authHandler.Refresh))
}
