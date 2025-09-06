package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
)

type AuthService interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	RegisterUser(ctx context.Context, email string, hashedPassword []byte) (int, error)
	SaveRefreshToken(ctx context.Context, userID int, token string) error
	RemoveRefreshToken(ctx context.Context, refreshToken string) error
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) AuthService {
	return &authService{
		userRepo,
		tokenRepo,
	}
}

func (s *authService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *authService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *authService) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.userRepo.CheckUserExistsByEmail(ctx, email)
}

func (s *authService) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (int, error) {
	return s.userRepo.RegisterUser(ctx, email, hashedPassword)
}

func (s *authService) SaveRefreshToken(ctx context.Context, userID int, token string) error {
	return s.tokenRepo.SaveRefreshToken(ctx, userID, token)
}

func (s *authService) RemoveRefreshToken(ctx context.Context, refreshToken string) error {
	return s.tokenRepo.RemoveRefreshTokenFromDB(ctx, refreshToken)
}
