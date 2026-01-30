package storage

import (
	"auth-service/internal/config"
	"auth-service/internal/logger"
	"context"
	"fmt"

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

func (p *PostgresGORM) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}

	logger.Log.Info().Msg("Closing postgres connection...")
	return sqlDB.Close()
}
