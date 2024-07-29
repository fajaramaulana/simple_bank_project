package controller

import (
	"context"
	"errors"

	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/handler/token"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/handler/validate"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/middleware"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/service"
	"github.com/fajaramaulana/simple_bank_project/pb"
)

type UserController struct {
	pb.UnimplementedSimpleBankServer
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) CreateUser(ctx context.Context, req *pb.CreateUserRequest, payload *token.Payload) (*pb.CreateUserRespose, error) {
	violations := validate.ValidateCreateUserRequest(req)

	if violations != nil {
		return nil, helper.InvalidArgumentError(violations)
	}

	res, err := c.userService.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *UserController) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, payload *token.Payload) (*pb.UpdateUserResponse, error) {
	// convert to string because the payload.UserUUID is a byte
	uuidUserReq, err := helper.ConvertStringToUUID(req.GetUserUuid())
	if err != nil {
		return nil, err
	}

	err = middleware.CheckRole(payload, "admin")
	if err != nil {
		return nil, helper.UnauthenticatedError(err)
	}

	if payload.UserUUID != uuidUserReq {
		return nil, helper.UnauthenticatedError(errors.New("access denied: user uuid is not match"))
	}
	violations := validate.ValidateUpdateUserRequest(req)
	if violations != nil {
		return nil, helper.InvalidArgumentError(violations)
	}
	res, err := c.userService.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
