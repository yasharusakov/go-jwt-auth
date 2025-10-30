package storage

import (
	"auth-service/internal/config"
	"context"
	"fmt"

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

	//err = createTables(ctx, pool)
	//if err != nil {
	//	return nil, err
	//}

	return pool, nil
}

//func createTables(ctx context.Context, pool *pgxpool.Pool) error {
//	_, err := pool.Exec(ctx, `
//			CREATE TABLE IF NOT EXISTS refresh_tokens (
//			id SERIAL PRIMARY KEY,
//			user_id INTEGER UNIQUE NOT NULL,
//			token TEXT UNIQUE NOT NULL,
//			expires_at TIMESTAMP NOT NULL,
//			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//		);
//	`)
//
//	if err != nil {
//		return fmt.Errorf("failed to create refresh_tokens table: %w", err)
//	}
//
//	log.Println("Database tables created or already exist.")
//
//	return nil
//}
