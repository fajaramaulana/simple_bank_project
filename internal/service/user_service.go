package service

import (
	"context"
	"errors"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/response"
	"github.com/fajaramaulana/simple_bank_project/util"
)

type UserService struct {
	db db.Store
}

func NewUserService(db db.Store) *UserService {
	return &UserService{
		db: db,
	}
}

func (u *UserService) CreateUser(ctx context.Context, request *request.CreateUserRequest) (response.UserResponseCreate, error) {
	// check if user already exists
	user, err := u.db.GetUserByEmail(ctx, request.Email)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return response.UserResponseCreate{
				Email: user.Email,
			}, err
		}
	}

	if len(user.Email) > 0 {
		return response.UserResponseCreate{
			Email: user.Email,
		}, errors.New("email already exists")
	}

	checkUsername, err := u.db.GetUserByUsername(ctx, request.Username)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return response.UserResponseCreate{
				Email: checkUsername.Email,
			}, err
		}
	}

	if len(checkUsername.Email) > 0 {
		return response.UserResponseCreate{
			Email: checkUsername.Email,
		}, errors.New("username already exists")
	}

	// create user
	hashPass, err := util.MakePasswordBcrypt(request.Password)
	if err != nil {
		return response.UserResponseCreate{}, err
	}
	request.Password = hashPass
	userCreate, err := u.db.CreateUserWithAccountTx(ctx, *request)
	if err != nil {
		return response.UserResponseCreate{}, err
	}

	return response.UserResponseCreate{
		UserUUID: userCreate.User.UserUUID,
		Username: userCreate.User.Username,
		FullName: userCreate.User.FullName,
		Email:    userCreate.User.Email,
		Account: response.AccountResponseSimple{
			AccountUUID: userCreate.Account.AccountUUID,
			Owner:       userCreate.Account.Owner,
			Currency:    userCreate.Account.Currency,
			Balance:     userCreate.Account.Balance,
		},
	}, nil
}
