package client

import (
	"context"

	userpb "tg-bot/internal/genproto/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserService interface {
	GetUsersCount(ctx context.Context) (int32, error)
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

func (u *UserClient) GetUsersCount(ctx context.Context) (int32, error) {
	resp, err := u.client.GetUsersCount(ctx, &userpb.GetUsersCountRequest{})
	if err != nil {
		return 0, err
	}

	return resp.Count, nil
}
