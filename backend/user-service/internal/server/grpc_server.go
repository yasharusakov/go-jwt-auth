package server

import (
	"net"
	userpb "user-service/internal/genproto/user"
	grpcHandler "user-service/internal/handler/grpc"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	GRPCServer *grpc.Server
}

func (s *GRPCServer) Run(port string, handlers *grpcHandler.UserHandler) error {
	s.GRPCServer = grpc.NewServer()
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	userpb.RegisterUserServiceServer(s.GRPCServer, handlers)
	return s.GRPCServer.Serve(lis)
}
