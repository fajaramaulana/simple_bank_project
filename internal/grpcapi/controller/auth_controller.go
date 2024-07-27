package controller

import (
	"context"

	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/service"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/shared"
	"github.com/fajaramaulana/simple_bank_project/pb"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	metaData := shared.ExtractMetadata(ctx)
	res, err := c.authService.LoginUser(ctx, req, metaData)
	if err != nil {
		return nil, err
	}

	return res, nil
}
