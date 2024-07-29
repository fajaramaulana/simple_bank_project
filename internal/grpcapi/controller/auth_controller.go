package controller

import (
	"context"

	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/handler/validate"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/service"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/shared"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validate.ValidateLoginUserRequest(req)
	if violations != nil {
		log.Err(helper.InvalidArgumentError(violations)).Msg("LoginUserRequest is invalid")
		return nil, helper.InvalidArgumentError(violations)
	}
	metaData := shared.ExtractMetadata(ctx)
	res, err := c.authService.LoginUser(ctx, req, metaData)
	if err != nil {
		log.Err(err).Msg("Failed to login user")
		return nil, err
	}

	return res, nil
}
