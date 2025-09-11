package service

import "user-service/internal/repository"

type Services struct {
	User UserService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		User: NewUserService(repos.User),
	}
}
