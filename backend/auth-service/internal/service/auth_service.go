package service

import (
	"auth-service/internal/client"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"context"
)

type AuthService interface {
	client.UserService
	SaveRefreshToken(ctx context.Context, userID string, token string) error
	RemoveRefreshToken(ctx context.Context, refreshToken string) error
}

type authService struct {
	userClient client.UserService
	tokenRepo  repository.TokenRepository
}

func NewAuthService(userClient client.UserService, tokenRepo repository.TokenRepository) AuthService {
	return &authService{
		userClient,
		tokenRepo,
	}
}

func (a authService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return a.userClient.GetUserByEmail(ctx, email)
}

func (a authService) GetUserByID(ctx context.Context, id string) (*model.UserWithoutPassword, error) {
	return a.userClient.GetUserByID(ctx, id)
}

func (a authService) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return a.userClient.CheckUserExistsByEmail(ctx, email)
}

func (a authService) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	return a.userClient.RegisterUser(ctx, email, hashedPassword)
}

func (a authService) SaveRefreshToken(ctx context.Context, userID string, token string) error {
	return a.tokenRepo.SaveRefreshToken(ctx, userID, token)
}

func (a authService) RemoveRefreshToken(ctx context.Context, refreshToken string) error {
	return a.tokenRepo.RemoveRefreshTokenFromDB(ctx, refreshToken)
}
