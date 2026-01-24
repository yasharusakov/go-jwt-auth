package repository

import (
	"context"
	"errors"
	"fmt"
	"user-service/internal/entity"

	"gorm.io/gorm"
)

type UserGormRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error)
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
}

type userGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) UserGormRepository {
	return &userGormRepository{db: db}
}

func (u *userGormRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	result := u.db.WithContext(ctx).Where("email = ?", email).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", result.Error)
		}
		return nil, fmt.Errorf("error getting user by email: %w", result.Error)
	}
	return user, nil
}

func (u *userGormRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}
	result := u.db.WithContext(ctx).Select("id, email").Where("id = ?", id).First(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", result.Error)
		}
		return nil, fmt.Errorf("error getting user by id: %w", result.Error)
	}
	return user, nil
}

func (u *userGormRepository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := u.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("check user exists by email error: %w", result.Error)
	}

	return count > 0, nil
}

func (u *userGormRepository) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	user := &entity.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	result := u.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return "", fmt.Errorf("error registering user: %w", result.Error)
	}

	return user.ID, nil
}

func (u *userGormRepository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	users := []*entity.User{}
	result := u.db.WithContext(ctx).Select("id, email").Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting all users: %w", result.Error)
	}
	return users, nil
}
