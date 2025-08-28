package repository

import (
	"context"
	"fmt"
	"server/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	RegisterUser(ctx context.Context, email string, hashedPassword []byte) (int, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	row := r.db.QueryRow(ctx, "SELECT id, email, password FROM users WHERE email=$1", email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	user := &model.User{}
	row := r.db.QueryRow(ctx, "SELECT id, email, password FROM users WHERE id=$1", id)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (int, error) {
	var userID int
	err := r.db.QueryRow(
		ctx,
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		email, hashedPassword,
	).Scan(&userID)
	return userID, err
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Query(ctx, "SELECT id, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
