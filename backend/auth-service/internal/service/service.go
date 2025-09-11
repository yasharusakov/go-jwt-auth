package service

import (
	"auth-service/internal/client"
	"auth-service/internal/repository"
)

type Services struct {
	Auth AuthService
}

func NewServices(userClient client.UserService, repos *repository.Repositories) *Services {
	return &Services{
		Auth: NewAuthService(userClient, repos.Token),
	}
}
