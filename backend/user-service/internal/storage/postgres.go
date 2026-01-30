package storage

import (
	"context"
	"fmt"
	"user-service/internal/config"
	"user-service/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type PostgresGORM struct {
	DB *gorm.DB
}

func NewPostgresGORM(ctx context.Context, cfg config.PostgresConfig) (*PostgresGORM, error) {
	postgresURI := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
		cfg.PostgresSSLMode,
	)

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database with GORM: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresGORM{DB: db}, nil
}

func (p *PostgresGORM) Close() {
	sqlDB, err := p.DB.DB()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to close database connection.")
	}

	logger.Log.Info().Msg("Closing postgres connection...")

	if err := sqlDB.Close(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Error closing postgres connection.")
	}
	logger.Log.Info().Msg("Database connection closed.")
}
