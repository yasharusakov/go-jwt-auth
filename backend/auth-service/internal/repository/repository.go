package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	Token TokenRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		Token: NewTokenRepository(db),
	}
}
