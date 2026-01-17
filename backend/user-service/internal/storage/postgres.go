package storage

import (
	"context"
	"fmt"
	"user-service/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	postgresUri := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
		cfg.PostgresSSLMode,
	)

	pool, err := pgxpool.New(ctx, postgresUri)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping storage: %w", err)
	}

	return pool, nil
}
