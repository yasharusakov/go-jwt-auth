package grpcClient

import (
	userpb "auth-service/internal/genproto/user/v1"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*userpb.GetUserByEmailResponse, error)
	GetUserByID(ctx context.Context, id string) (*userpb.GetUserByIDResponse, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (*userpb.CheckUserExistsByEmailResponse, error)
	RegisterUser(ctx context.Context, email string, hashedPassword []byte) (*userpb.RegisterUserResponse, error)
	Ping(ctx context.Context) error
}

type UserClient struct {
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := userpb.NewUserServiceClient(conn)

	return &UserClient{
		client: c,
		conn:   conn,
	}, nil
}

func (u *UserClient) Ping(ctx context.Context) error {
	state := u.conn.GetState()
	if state != connectivity.Ready && state != connectivity.Idle {
		return fmt.Errorf("grpc connection is not ready: %s", state.String())
	}
	return nil
}

func (u *UserClient) Close() {
	_ = u.conn.Close()
}

func (u *UserClient) GetUserByEmail(ctx context.Context, email string) (*userpb.GetUserByEmailResponse, error) {
	resp, err := u.client.GetUserByEmail(ctx, &userpb.GetUserByEmailRequest{Email: email})
	if err != nil {
		return nil, err
	}
	if resp.User == nil {
		return nil, fmt.Errorf("user not found")
	}
	return &userpb.GetUserByEmailResponse{
		User: resp.User,
	}, nil
}

func (u *UserClient) GetUserByID(ctx context.Context, id string) (*userpb.GetUserByIDResponse, error) {
	resp, err := u.client.GetUserByID(ctx, &userpb.GetUserByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}
	if resp.User == nil {
		return nil, fmt.Errorf("user not found")
	}
	return &userpb.GetUserByIDResponse{
		User: resp.User,
	}, nil
}

func (u *UserClient) CheckUserExistsByEmail(ctx context.Context, email string) (*userpb.CheckUserExistsByEmailResponse, error) {
	resp, err := u.client.CheckUserExistsByEmail(ctx, &userpb.CheckUserExistsByEmailRequest{Email: email})
	if err != nil {
		return nil, err
	}
	return &userpb.CheckUserExistsByEmailResponse{
		Exists: resp.Exists,
	}, nil
}

func (u *UserClient) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (*userpb.RegisterUserResponse, error) {
	resp, err := u.client.RegisterUser(ctx, &userpb.RegisterUserRequest{
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}
	return &userpb.RegisterUserResponse{
		Id: resp.Id,
	}, nil
}
