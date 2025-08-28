package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.GetAllUsers(ctx)
}
