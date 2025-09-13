package handler

import "auth-service/internal/service"

type Handlers struct {
	Auth AuthHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(services.Auth),
	}
}
