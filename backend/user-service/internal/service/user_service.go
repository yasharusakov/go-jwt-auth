package service

import (
	"context"
	"user-service/internal/repository"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserWithoutPassword struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*UserWithoutPassword, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error)
	GetAllUsers(ctx context.Context) ([]*UserWithoutPassword, error)
}

type userService struct {
	repo repository.UserGormRepository
}

func NewUserService(repo repository.UserGormRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	result, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       result.ID,
		Email:    result.Email,
		Password: result.Password,
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*UserWithoutPassword, error) {
	result, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserWithoutPassword{
		ID:    result.ID,
		Email: result.Email,
	}, nil
}

func (s *userService) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.repo.CheckUserExistsByEmail(ctx, email)
}

func (s *userService) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	return s.repo.RegisterUser(ctx, email, hashedPassword)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*UserWithoutPassword, error) {
	result, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var users []*UserWithoutPassword
	for _, u := range result {
		users = append(users, &UserWithoutPassword{
			ID:    u.ID,
			Email: u.Email,
		})
	}

	return users, nil
}
