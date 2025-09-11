package router

import (
	"net/http"
	"user-service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(handlers *handler.Handlers) http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/api/user/get-by-email/{email}", handlers.User.GetUserByEmail).Methods("GET")
	m.HandleFunc("/api/user/get-by-id/{id}", handlers.User.GetUserByID).Methods("GET")
	m.HandleFunc("/api/user/check-by-email/{email}", handlers.User.CheckUserExistsByEmail).Methods("GET")
	m.HandleFunc("/api/user/register", handlers.User.RegisterUser).Methods("POST")
	m.HandleFunc("/api/user/get-all", handlers.User.GetAllUsers).Methods("GET")

	return m
}
