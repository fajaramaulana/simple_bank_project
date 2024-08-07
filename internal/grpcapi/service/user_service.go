package service

import (
	"context"
	"fmt"
	"time"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	db          db.Store
	config      util.Config
	redisClient *redis.Client
}

func NewUserService(db db.Store, config util.Config, redisClient *redis.Client) *UserService {
	return &UserService{db: db, config: config, redisClient: redisClient}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserRespose, error) {
	user, err := s.db.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		fmt.Printf("%# v\n", err.Error())
		if err.Error() != "no rows in result set" {
			// return nil, err
			return nil, status.Errorf(codes.Internal, "failed to get user by email: %v", err)
		}
	}

	if len(user.Email) > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}

	checkUsername, err := s.db.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		fmt.Printf("%# v\n", err.Error())
		if err.Error() != "no rows in result set" {
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

	// make verification code and save to redis

	verificationCode := uuid.New().String()

	uuidUser, err := helper.ConvertStringToUUID(userCreate.User.UserUUID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_uuid: %v", err)
	}

	argVerCode := db.UpdateUserVerificationEmailParams{
		UserUuid: uuidUser,
		VerificationEmailCode: pgtype.Text{
			String: verificationCode,
			Valid:  true,
		},
		VerificationEmailExpiredAt: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Minute * 15),
			Valid: true,
		},
	}

	result, err := s.db.UpdateUserVerificationEmail(ctx, argVerCode)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user verification email: %v", err)
	}

	// insert to redis
	keyRedis := "verification_email:users:" + userCreate.User.UserUUID
	valueRedis := map[string]interface{}{
		"verification_code": verificationCode,
		"expired_at":        result.VerificationEmailExpiredAt.Time,
		"email":             userCreate.User.Email,
	}
	_, err = s.redisClient.HSet(ctx, keyRedis, valueRedis).Result()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set verification code in Redis: %v", err)
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
		if err.Error() != "no rows in result set" {
			return nil, status.Errorf(codes.Internal, "failed to get user by email: %v", err)
		}
	}

	if len(checkEmail.Email) > 0 && checkEmail.UserUuid != uuidUser {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}
	arg := db.UpdateUserParams{
		UserUuid: uuidUser,
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
	}

	if req.Password != nil {
		hashPass, err := util.MakePasswordBcrypt(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
		}

		arg.HashedPassword = pgtype.Text{
			String: hashPass,
			Valid:  true,
		}
		arg.PasswordChangedAt = pgtype.Timestamptz{
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
			Balance:     account.Balance.Int.String(),
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
