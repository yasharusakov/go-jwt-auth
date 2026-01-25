package storage

import (
	"auth-service/internal/config"
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//func NewPostgres(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
//	postgresUri := fmt.Sprintf(
//		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
//		cfg.PostgresUser,
//		cfg.PostgresPassword,
//		cfg.PostgresHost,
//		cfg.PostgresPort,
//		cfg.PostgresDB,
//		cfg.PostgresSSLMode,
//	)
//
//	pool, err := pgxpool.New(ctx, postgresUri)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create connection pool: %w", err)
//	}
//
//	err = pool.Ping(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("failed to ping storage: %w", err)
//	}
//
//	return pool, nil
//}

type PostgresGORM struct {
	DB *gorm.DB
}

func NewPostgresGORM(ctx context.Context, cfg config.PostgresConfig) (*PostgresGORM, error) {
	postgresUri := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
		cfg.PostgresSSLMode,
	)

	db, err := gorm.Open(postgres.Open(postgresUri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
	return sqlDB.Close()
}
