package router

import (
	"net/http"
	"server/internal/handler"
	"server/internal/middleware"
)

func RegisterUserRoutes(mux *http.ServeMux, userHandler handler.UserHandler) {
	mux.HandleFunc("/api/users",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(userHandler.GetAllUsers),
		),
	)
}
