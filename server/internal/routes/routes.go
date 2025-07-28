package routes

import (
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	RegisterAuthRoutes(mux)
	RegisterUserRoutes(mux)
	return mux
}
