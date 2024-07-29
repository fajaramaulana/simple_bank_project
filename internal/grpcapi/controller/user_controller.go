package controller

import (
	"context"
	"fmt"

	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/handler/validate"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
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

func (c *UserController) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserRespose, error) {
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

func (c *UserController) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	fmt.Printf("%# v\n", req)
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
