package service

import (
	grpcClient "auth-service/internal/client/grpc"
	userpb "auth-service/internal/genproto/user"
	"auth-service/internal/repository"
	"context"
)

type AuthService interface {
	//httpClient.UserService
	grpcClient.UserService
	repository.TokenRepository
}

type authService struct {
	//httpUserClient httpClient.UserService
	grpcUserClient grpcClient.UserService
	tokenRepo      repository.TokenRepository
}

// HTTP client version
//func NewAuthService(httpUserClient httpClient.UserService, tokenRepo repository.TokenRepository) AuthService {
//	return &authService{
//		httpUserClient,
//		tokenRepo,
//	}
//}

func NewAuthService(grpcUserClient grpcClient.UserService, tokenRepo repository.TokenRepository) AuthService {
	return &authService{
		grpcUserClient,
		tokenRepo,
	}
}

func (s authService) GetUserByEmail(ctx context.Context, email string) (*userpb.GetUserByEmailResponse, error) {
	return s.grpcUserClient.GetUserByEmail(ctx, email)
}

func (s authService) GetUserByID(ctx context.Context, id string) (*userpb.GetUserByIDResponse, error) {
	return s.grpcUserClient.GetUserByID(ctx, id)
}

func (s authService) CheckUserExistsByEmail(ctx context.Context, email string) (*userpb.CheckUserExistsByEmailResponse, error) {
	return s.grpcUserClient.CheckUserExistsByEmail(ctx, email)
}

func (s authService) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (*userpb.RegisterUserResponse, error) {
	return s.grpcUserClient.RegisterUser(ctx, email, hashedPassword)
}

func (s authService) SaveRefreshToken(ctx context.Context, userID string, token string) error {
	return s.tokenRepo.SaveRefreshToken(ctx, userID, token)
}

func (s authService) RemoveRefreshToken(ctx context.Context, refreshToken string) error {
	return s.tokenRepo.RemoveRefreshToken(ctx, refreshToken)
}
