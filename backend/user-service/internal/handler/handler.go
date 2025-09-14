package handler

import (
	"user-service/internal/handler/http"
	"user-service/internal/service"
)

type Handlers struct {
	User http.UserHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		User: http.NewUserHandler(services.User),
	}
}
