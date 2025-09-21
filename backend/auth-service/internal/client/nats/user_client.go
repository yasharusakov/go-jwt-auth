package natsClient

import (
	"auth-service/internal/model"
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.UserWithoutPassword, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error)
}

type UserClient struct {
	nc *nats.Conn
}

func NewUserClient(nc *nats.Conn) UserService {
	return &UserClient{nc}
}

func (u *UserClient) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	subject := "user.getByEmail"
	data, _ := json.Marshal(map[string]string{"email": email})

	msg, err := u.nc.RequestWithContext(ctx, subject, data)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal(msg.Data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserClient) GetUserByID(ctx context.Context, id string) (*model.UserWithoutPassword, error) {
	subject := "user.getById"
	data, _ := json.Marshal(map[string]string{"id": id})

	msg, err := u.nc.RequestWithContext(ctx, subject, data)
	if err != nil {
		return nil, err
	}

	var user model.UserWithoutPassword
	if err := json.Unmarshal(msg.Data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserClient) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	subject := "user.checkByEmail"
	data, _ := json.Marshal(map[string]string{"email": email})

	msg, err := u.nc.RequestWithContext(ctx, subject, data)
	if err != nil {
		return false, err
	}

	var result struct {
		Exists bool `json:"exists"`
	}
	if err = json.Unmarshal(msg.Data, &result); err != nil {
		return false, err
	}

	return result.Exists, nil
}

func (u *UserClient) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	subject := "user.register"
	payload := map[string]interface{}{
		"email":    email,
		"password": hashedPassword,
	}
	data, _ := json.Marshal(payload)

	msg, err := u.nc.RequestWithContext(ctx, subject, data)
	if err != nil {
		return "", err
	}

	var result struct {
		ID string `json:"id"`
	}
	if err = json.Unmarshal(msg.Data, &result); err != nil {
		return "", err
	}

	return result.ID, nil
}
