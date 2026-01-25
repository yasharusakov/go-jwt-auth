package repository

import (
	"auth-service/internal/config"
	"auth-service/internal/entity"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error
	RemoveRefreshToken(ctx context.Context, refreshToken string) error
	IsRefreshTokenExists(ctx context.Context, refreshToken string) (bool, error)
}

type tokenRepository struct {
	db  *gorm.DB
	cfg config.Config
}

func NewTokenRepository(db *gorm.DB, cfg config.Config) TokenRepository {
	return &tokenRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *tokenRepository) SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	//_, err := r.db.Exec(ctx,
	//	"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
	//	userID, refreshToken, time.Now().Add(r.cfg.JWT.JWTRefreshTokenExp),
	//)
	//if err != nil {
	//	return fmt.Errorf("failed to save refresh token: %w", err)
	//}
	token := entity.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(r.cfg.JWT.JWTRefreshTokenExp),
	}

	result := r.db.WithContext(ctx).Create(&token)
	if result.Error != nil {
		return fmt.Errorf("failed to save refresh token: %w", result.Error)
	}

	return nil
}

func (r *tokenRepository) RemoveRefreshToken(ctx context.Context, refreshToken string) error {
	//_, err := r.db.Exec(ctx, "DELETE FROM refresh_tokens WHERE token = $1", refreshToken)
	//if err != nil {
	//	return fmt.Errorf("failed to remove refresh token: %w", err)
	//}

	result := r.db.WithContext(ctx).Where("token = ?", refreshToken).Delete(&entity.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove refresh token: %w", result.Error)
	}

	return nil
}

func (r *tokenRepository) IsRefreshTokenExists(ctx context.Context, refreshToken string) (bool, error) {
	//var exists bool
	//err := r.db.QueryRow(ctx,
	//	"SELECT EXISTS(SELECT 1 FROM refresh_tokens WHERE token = $1)",
	//	refreshToken,
	//).Scan(&exists)
	//if err != nil {
	//	return false, fmt.Errorf("failed to check refresh token existence: %w", err)
	//}

	var count int64
	result := r.db.WithContext(ctx).Model(&entity.RefreshToken{}).Where("token = ?", refreshToken).Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check refresh token existence: %w", result.Error)
	}

	return count > 0, nil
}
