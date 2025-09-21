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
	SaveRefreshToken(ctx context.Context, userID string, token string) error
	RemoveRefreshToken(ctx context.Context, refreshToken string) error
}

type authService struct {
	//httpUserClient httpClient.UserService
	grpcUserClient grpcClient.UserService
	tokenRepo      repository.TokenRepository
}

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

func (a authService) GetUserByEmail(ctx context.Context, email string) (*userpb.GetUserByEmailResponse, error) {
	return a.grpcUserClient.GetUserByEmail(ctx, email)
}

func (a authService) GetUserByID(ctx context.Context, id string) (*userpb.GetUserByIDResponse, error) {
	return a.grpcUserClient.GetUserByID(ctx, id)
}

func (a authService) CheckUserExistsByEmail(ctx context.Context, email string) (*userpb.CheckUserExistsByEmailResponse, error) {
	return a.grpcUserClient.CheckUserExistsByEmail(ctx, email)
}

func (a authService) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (*userpb.RegisterUserResponse, error) {
	return a.grpcUserClient.RegisterUser(ctx, email, hashedPassword)
}

func (a authService) SaveRefreshToken(ctx context.Context, userID string, token string) error {
	return a.tokenRepo.SaveRefreshToken(ctx, userID, token)
}

func (a authService) RemoveRefreshToken(ctx context.Context, refreshToken string) error {
	return a.tokenRepo.RemoveRefreshTokenFromDB(ctx, refreshToken)
}
