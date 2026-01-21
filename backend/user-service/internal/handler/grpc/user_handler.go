package grpcHandler

import (
	"context"
	userpb "user-service/internal/genproto/user/v1"
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

// We do not use gRPC method to get all users anymore
// due to http endpoint implementation
//func (h *UserHandler) GetAllUsers(ctx context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error) {
//	users, err := h.userService.GetAllUsers(ctx)
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "error getting users: %v", err)
//	}
//
//	newUsers := make([]*userpb.User, len(users))
//	for i, user := range users {
//		newUsers[i] = &userpb.User{
//			Id:    user.ID,
//			Email: user.Email,
//		}
//	}
//
//	return &userpb.GetAllUsersResponse{
//		Users: newUsers,
//	}, nil
//}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserByEmailResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user by email: %v", err)
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
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
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := h.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user by id: %v", err)
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &userpb.GetUserByIDResponse{
		User: &userpb.User{
			Id:    user.ID,
			Email: user.Email,
		},
	}, nil
}

func (h *UserHandler) CheckUserExistsByEmail(ctx context.Context, req *userpb.CheckUserExistsByEmailRequest) (*userpb.CheckUserExistsByEmailResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	exists, err := h.userService.CheckUserExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check user existence: %v", err)
	}
	return &userpb.CheckUserExistsByEmailResponse{Exists: exists}, nil
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.RegisterUserResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	userID, err := h.userService.RegisterUser(ctx, req.Email, []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &userpb.RegisterUserResponse{
		Id: userID,
	}, nil
}
