package handler

import (
	"server/internal/service"
)

type Handlers struct {
	Auth AuthHandler
	User UserHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(services.Auth),
		User: NewUserHandler(services.User),
	}
}
