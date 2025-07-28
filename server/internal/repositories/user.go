package repositories

import (
	"context"
	"server/internal/database/postgresql"
	"server/internal/models"
)

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	row := postgresql.Pool.QueryRow(ctx, "SELECT id, email, password FROM users WHERE email=$1", email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	row := postgresql.Pool.QueryRow(ctx, "SELECT id, email, password FROM users WHERE id=$1", id)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := postgresql.Pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func RegisterUser(ctx context.Context, email string, hashedPassword []byte) (int, error) {
	var userID int
	err := postgresql.Pool.QueryRow(
		ctx,
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		email, hashedPassword,
	).Scan(&userID)
	return userID, err
}
