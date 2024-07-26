package grpcapi

import (
	"context"
	"time"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/fajaramaulana/simple_bank_project/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserRespose, error) {
	user, err := s.db.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			// return nil, err
			return nil, status.Errorf(codes.Internal, "failed to get user by email: %v", err)
		}
	}

	if len(user.Email) > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}

	checkUsername, err := s.db.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, status.Errorf(codes.Internal, "failed to get user by username: %v", err)
		}
	}

	if len(checkUsername.Email) > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	}

	hashPass, err := util.MakePasswordBcrypt(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	arg := request.CreateUserRequest{
		Username: req.GetUsername(),
		Password: hashPass,
		Email:    req.GetEmail(),
		FullName: req.GetFullName(),
		Currency: req.GetCurrency(),
	}
	userCreate, err := s.db.CreateUserWithAccountTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	res := &pb.CreateUserRespose{
		User: &pb.User{
			UserUuid:  userCreate.User.UserUUID,
			Username:  userCreate.User.Username,
			FullName:  userCreate.User.FullName,
			Email:     userCreate.User.Email,
			CreatedAt: timestamppb.New(time.Now()),
		},
		Account: &pb.Account{
			AccountUuid: userCreate.Account.AccountUUID.String(),
			Owner:       userCreate.Account.Owner,
			Currency:    userCreate.Account.Currency,
			Balance:     userCreate.Account.Balance,
			CreatedAt:   timestamppb.New(time.Now()),
		},
	}

	return res, nil
}
