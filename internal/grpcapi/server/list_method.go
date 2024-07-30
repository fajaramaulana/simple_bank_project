package server

import (
	"context"

	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/middleware"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/rs/zerolog/log"
)

// Implement gRPC methods using the controllers
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserRespose, error) {
	// s.tokenMaker
	payload, err := middleware.AuthMiddleware(ctx, s.tokenMaker)
	if err != nil {
		log.Error().Err(err).Msg("Failed to authenticate")
		return nil, err
	}
	return s.userController.CreateUser(ctx, req, payload)
}
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	payload, err := middleware.AuthMiddleware(ctx, s.tokenMaker)
	if err != nil {
		log.Error().Err(err).Msg("Failed to authenticate")
		return nil, err
	}

	return s.userController.UpdateUser(ctx, req, payload)
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return s.authController.LoginUser(ctx, req)
}

func (s *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	return s.authController.VerifyEmail(ctx, req)
}
