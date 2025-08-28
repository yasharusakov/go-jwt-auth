package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User  UserRepository
	Token TokenRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:  NewUserRepository(db),
		Token: NewTokenRepository(db),
	}
}
