package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"server/internal/config"
)

var (
	Pool *pgxpool.Pool
)

func createTables(ctx context.Context) error {
	_, err := Pool.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
	`)

	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = Pool.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token TEXT UNIQUE NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		return fmt.Errorf("failed to create refresh_tokens table: %w", err)
	}

	log.Println("Database tables created or already exist.")

	return nil
}

func Connect(config config.PostgresqlConfig) error {
	var err error
	ctx := context.Background()

	Pool, err = pgxpool.New(ctx, config.PostgresqlUri)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	err = Pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping the database: %w", err)
	}

	err = createTables(ctx)
	if err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	log.Println("Connected to postgres database successfully.")

	return nil
}
