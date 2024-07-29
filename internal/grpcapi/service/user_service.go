package service

import (
	"context"
	"database/sql"
	"time"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/fajaramaulana/simple_bank_project/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	db     db.Store
	config util.Config
}

func NewUserService(db db.Store, config util.Config) *UserService {
	return &UserService{db: db, config: config}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserRespose, error) {
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

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	uuidUser, err := helper.ConvertStringToUUID(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_uuid: %v", err)
	}

	_, err = s.db.GetUserByUserUUID(ctx, uuidUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user by user_uuid: %v", err)
	}

	checkEmail, err := s.db.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, status.Errorf(codes.Internal, "failed to get user by email: %v", err)
		}
	}

	if len(checkEmail.Email) > 0 && checkEmail.UserUuid != uuidUser {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}
	arg := db.UpdateUserParams{
		UserUuid: uuidUser,
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
	}

	if req.Password != nil {
		hashPass, err := util.MakePasswordBcrypt(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
		}

		arg.HashedPassword = sql.NullString{
			String: hashPass,
			Valid:  true,
		}
		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	userUpdate, err := s.db.UpdateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	// list acocunt by user_uuid
	accounts, err := s.db.GetAccountByUserUUIDMany(ctx, userUpdate.UserUuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get account by user_uuid: %v", err)
	}

	if len(accounts) == 0 {
		return nil, status.Errorf(codes.Internal, "account not found")
	}

	// convert accounts to []*Account
	accountsRes := make([]*pb.Account, 0)
	for _, account := range accounts {
		accountsRes = append(accountsRes, &pb.Account{
			AccountUuid: account.AccountUuid.String(),
			Owner:       account.Owner,
			Currency:    account.Currency,
			Balance:     account.Balance,
			CreatedAt:   timestamppb.New(account.CreatedAt),
		})
	}

	res := &pb.UpdateUserResponse{
		UserUuid: userUpdate.UserUuid.String(),
		Username: userUpdate.Username,
		Email:    userUpdate.Email,
		FullName: userUpdate.FullName,
		Account:  accountsRes,
	}

	return res, nil
}
