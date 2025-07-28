package routes

import (
	"net/http"
	"server/internal/handlers/user"
	"server/internal/middlewares"
)

func RegisterUserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", middlewares.CORSMiddleware(middlewares.AuthMiddleware(user.GetUsers)))
}
