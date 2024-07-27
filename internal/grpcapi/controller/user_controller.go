package controller

import (
	"context"

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
	res, err := c.userService.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
