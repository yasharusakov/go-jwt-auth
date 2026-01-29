package service

import (
	grpcClient "auth-service/internal/client/grpc"
	"auth-service/internal/config"
	"auth-service/internal/logger"
	"auth-service/internal/repository"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	Email        string
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (*AuthResult, error)
	Login(ctx context.Context, email, password string) (*AuthResult, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthResult, error)
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

func (s *authService) Register(ctx context.Context, email, password string) (*AuthResult, error) {
	exists, err := s.grpcUserClient.CheckUserExistsByEmail(ctx, email)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error checking user existence")
		return nil, fmt.Errorf("internal server error")
	}
	if exists.Exists {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error hashing password")
		return nil, fmt.Errorf("internal server error")
	}

	userResp, err := s.grpcUserClient.RegisterUser(ctx, email, hashedPassword)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error gRPC register user")
		return nil, fmt.Errorf("internal server error")
	}

	return s.createSession(ctx, userResp.Id, email)
}

func (s *authService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	userResp, err := s.grpcUserClient.GetUserByEmail(ctx, email)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, fmt.Errorf("invalid email or password")
		}
		logger.Log.Error().Err(err).Msg("error gRPC get user by email")
		return nil, fmt.Errorf("internal server error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userResp.User.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return s.createSession(ctx, userResp.User.Id, email)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	claims, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error validate refresh token")
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	userID := claims.Subject

	// Verify user exists
	_, err = s.grpcUserClient.GetUserByID(ctx, userID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error gRPC get user by ID")
		return nil, fmt.Errorf("user not found")
	}

	isExists, err := s.tokenRepo.IsRefreshTokenExists(ctx, refreshToken)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error checking refresh token existence")
		return nil, fmt.Errorf("internal server error")
	}

	if !isExists {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	err = s.tokenRepo.RemoveRefreshToken(ctx, refreshToken)
	if err != nil {
		logger.Log.Error().Err(err).Msg("removing old refresh token failed")
		return nil, fmt.Errorf("internal server error")
	}

	return s.createSession(ctx, userID, refreshToken)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenRepo.RemoveRefreshToken(ctx, refreshToken)
}

func (s *authService) createSession(ctx context.Context, userID, email string) (*AuthResult, error) {
	accessToken, refreshToken, err := s.tokenManager.GenerateTokens(userID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error generating tokens")
		return nil, fmt.Errorf("internal server error")
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error saving refresh token")
		return nil, fmt.Errorf("internal server error")
	}

	return &AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userID,
		Email:        email,
	}, nil
}
