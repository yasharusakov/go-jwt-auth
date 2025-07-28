package repositories

import (
	"context"
	"server/internal/database/postgresql"
	"time"
)

func SaveRefreshToken(ctx context.Context, userID int, refreshToken string, expRefresh time.Duration) error {
	_, err := postgresql.Pool.Exec(ctx, "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, refreshToken, time.Now().Add(expRefresh))
	return err
}

func RemoveRefreshTokenFromDB(ctx context.Context, refreshToken string) error {
	_, err := postgresql.Pool.Exec(ctx, "DELETE FROM refresh_tokens WHERE token = $1", refreshToken)
	return err
}
