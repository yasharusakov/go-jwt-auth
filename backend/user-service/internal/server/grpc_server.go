package server

import (
	"context"
	"net"
	userpb "user-service/internal/genproto/user"
	grpcHandler "user-service/internal/handler/grpc"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	GRPCServer *grpc.Server
}

func (s *GRPCServer) Run(port string, handlers *grpcHandler.UserHandler) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s.GRPCServer = grpc.NewServer()
	userpb.RegisterUserServiceServer(s.GRPCServer, handlers)
	return s.GRPCServer.Serve(lis)
}

func (s *GRPCServer) Shutdown(ctx context.Context) {
	if s.GRPCServer == nil {
		return
	}

	done := make(chan struct{})
	go func() {
		s.GRPCServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		s.GRPCServer.Stop()
	}
}
