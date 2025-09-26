package grpcHandler

import (
	"context"
	userpb "user-service/internal/genproto/user"
	"user-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) GetAllUsers(ctx context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error) {
	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting users: %v", err)
	}

	newUsers := make([]*userpb.User, len(users))
	for i, user := range users {
		newUsers[i] = &userpb.User{
			Id:    user.ID,
			Email: user.Email,
		}
	}

	return &userpb.GetAllUsersResponse{
		Users: newUsers,
	}, nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserByEmailResponse, error) {
	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &userpb.GetUserByEmailResponse{}, nil
	}

	return &userpb.GetUserByEmailResponse{
		User: &userpb.User{
			Id:       user.ID,
			Email:    user.Email,
			Password: user.Password,
		},
	}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.GetUserByIDResponse, error) {
	user, err := h.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &userpb.GetUserByIDResponse{}, nil
	}

	return &userpb.GetUserByIDResponse{
		User: &userpb.User{
			Id:    user.ID,
			Email: user.Email,
		},
	}, nil
}

func (h *UserHandler) CheckUserExistsByEmail(ctx context.Context, req *userpb.CheckUserExistsByEmailRequest) (*userpb.CheckUserExistsByEmailResponse, error) {
	exists, err := h.userService.CheckUserExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &userpb.CheckUserExistsByEmailResponse{Exists: exists}, nil
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.RegisterUserResponse, error) {
	userID, err := h.userService.RegisterUser(ctx, req.Email, []byte(req.Password))
	if err != nil {
		return nil, err
	}

	return &userpb.RegisterUserResponse{
		Id: userID,
	}, nil
}
