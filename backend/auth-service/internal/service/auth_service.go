package service

import (
	"auth-service/internal/apperror"
	grpcClient "auth-service/internal/client/grpc"
	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/repository"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*dto.AuthResult, error)
	Login(ctx context.Context, email, password string) (*dto.AuthResult, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResult, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	grpcUserClient grpcClient.UserService
	tokenRepo      repository.TokenRepository
	tokenManager   TokenManager
	cfg            config.Config
}

func NewAuthService(
	grpcUserClient grpcClient.UserService,
	tokenRepo repository.TokenRepository,
	tokenManager TokenManager,
	cfg config.Config,
) AuthService {
	return &authService{
		grpcUserClient: grpcUserClient,
		tokenRepo:      tokenRepo,
		tokenManager:   tokenManager,
		cfg:            cfg,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (*dto.AuthResult, error) {
	exists, err := s.grpcUserClient.CheckUserExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check user existence: %w", err)
	}
	if exists.Exists {
		return nil, apperror.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	userResp, err := s.grpcUserClient.RegisterUser(ctx, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("gRPC register: %w", err)
	}

	return s.createSession(ctx, userResp.Id, email)
}

func (s *authService) Login(ctx context.Context, email, password string) (*dto.AuthResult, error) {
	userResp, err := s.grpcUserClient.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, apperror.ErrInvalidEmailOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(userResp.User.Password), []byte(password))
	if err != nil {
		return nil, apperror.ErrInvalidEmailOrPassword
	}

	return s.createSession(ctx, userResp.User.Id, email)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResult, error) {
	claims, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperror.ErrInvalidOrExpiredRefreshToken, err)
	}

	isExists, err := s.tokenRepo.IsRefreshTokenExists(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("checking refresh token existence: %w", err)
	}

	if !isExists {
		return nil, apperror.ErrInvalidOrExpiredRefreshToken
	}

	userID := claims.Subject

	err = s.tokenRepo.RemoveRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("removing old refresh token failed: %w", err)
	}

	accessToken, refreshToken, err := s.tokenManager.GenerateTokens(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperror.ErrGeneratingTokens, err)
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("saving refresh token failed: %w", err)
	}

	userData, err := s.grpcUserClient.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperror.ErrUserNotFound, err)
	}

	return &dto.RefreshResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:    userData.User.Id,
			Email: userData.User.Email,
		},
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenRepo.RemoveRefreshToken(ctx, refreshToken)
}

func (s *authService) createSession(ctx context.Context, userID, email string) (*dto.AuthResult, error) {
	accessToken, refreshToken, err := s.tokenManager.GenerateTokens(userID)
	if err != nil {
		return nil, fmt.Errorf("token generation failed: %w", err)
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("saving refresh token failed: %w", err)
	}

	return &dto.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:    userID,
			Email: email,
		},
	}, nil
}
