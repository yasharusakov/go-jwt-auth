package service

import "server/internal/repository"

type Services struct {
	Auth AuthService
	User UserService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Auth: NewAuthService(repos.User, repos.Token),
		User: NewUserService(repos.User),
	}
}
