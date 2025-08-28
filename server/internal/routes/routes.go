package routes

import (
	"net/http"
	"server/internal/handler"
)

func SetupRoutes(handlers *handler.Handlers) *http.ServeMux {
	mux := http.NewServeMux()
	RegisterAuthRoutes(mux, handlers.Auth)
	RegisterUserRoutes(mux, handlers.User)
	return mux
}
