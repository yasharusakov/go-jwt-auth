package repository

import (
	"auth-service/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error
	RemoveRefreshToken(ctx context.Context, refreshToken string) error
}

type tokenRepository struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewTokenRepository(db *pgxpool.Pool, cfg *config.Config) TokenRepository {
	return &tokenRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *tokenRepository) SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, refreshToken, time.Now().Add(r.cfg.JWT.JWTRefreshTokenExp),
	)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	return nil
}

func (r *tokenRepository) RemoveRefreshToken(ctx context.Context, refreshToken string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM refresh_tokens WHERE token = $1", refreshToken)
	if err != nil {
		return fmt.Errorf("failed to remove refresh token: %w", err)
	}
	return nil
}
