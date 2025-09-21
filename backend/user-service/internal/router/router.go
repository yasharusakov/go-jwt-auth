package router

import (
	"net/http"
	httpHandler "user-service/internal/handler/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(handlers httpHandler.Handlers) http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/api/user/get-by-email/{email}", handlers.GetUserByEmail).Methods("GET")
	m.HandleFunc("/api/user/get-by-id/{id}", handlers.GetUserByID).Methods("GET")
	m.HandleFunc("/api/user/check-by-email/{email}", handlers.CheckUserExistsByEmail).Methods("GET")
	m.HandleFunc("/api/user/register", handlers.RegisterUser).Methods("POST")
	m.HandleFunc("/api/user/get-all", handlers.GetAllUsers).Methods("GET")

	return m
}
