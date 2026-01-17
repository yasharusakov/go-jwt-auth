package httpHandler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/service"
)

type Handlers interface {
	GetAllUsers(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) Handlers {
	return &userHandler{service}
}

func (h *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
