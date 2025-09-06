package repository

import (
	"context"
	"server/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID int, refreshToken string) error
	RemoveRefreshTokenFromDB(ctx context.Context, refreshToken string) error
}

type tokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) SaveRefreshToken(ctx context.Context, userID int, refreshToken string) error {
	expRefresh := config.LoadConfig().JWT.JWTRefreshTokenExp
	_, err := r.db.Exec(ctx, "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, refreshToken, time.Now().Add(expRefresh))
	return err
}

func (r *tokenRepository) RemoveRefreshTokenFromDB(ctx context.Context, refreshToken string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM refresh_tokens WHERE token = $1", refreshToken)
	return err
}
