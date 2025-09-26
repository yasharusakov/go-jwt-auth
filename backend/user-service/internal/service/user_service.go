package service

import (
	"context"
	"user-service/internal/model"
	"user-service/internal/repository"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.UserWithoutPassword, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error)
	GetAllUsers(ctx context.Context) ([]model.UserWithoutPassword, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.UserWithoutPassword, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *userService) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.repo.CheckUserExistsByEmail(ctx, email)
}

func (s *userService) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	return s.repo.RegisterUser(ctx, email, hashedPassword)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]model.UserWithoutPassword, error) {
	return s.repo.GetAllUsers(ctx)
}
