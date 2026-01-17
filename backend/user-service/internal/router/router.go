package router

import (
	"net/http"
	httpHandler "user-service/internal/handler/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(handlers httpHandler.Handlers) http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/api/user/get-all", handlers.GetAllUsers).Methods("GET")

	return m
}
