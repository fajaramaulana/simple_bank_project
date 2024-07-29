package server

import (
	"context"

	"github.com/fajaramaulana/simple_bank_project/pb"
)

// Implement gRPC methods using the controllers
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserRespose, error) {
	// metaData := s.ExtractMetadata(ctx)
	return s.userController.CreateUser(ctx, req)
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return s.authController.LoginUser(ctx, req)
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return s.userController.UpdateUser(ctx, req)
}
