package service

import (
	"context"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/response"
	"github.com/google/uuid"
)

type AccountService struct {
	db db.Store
}

func NewAccountService(db db.Store) *AccountService {
	return &AccountService{
		db: db,
	}
}

func (a *AccountService) CreateAccount(ctx context.Context, request *request.CreateAccountRequest) (response.AccountResponseCreate, error) {
	// convert req.uuseruuid to uuid
	useruuid, err := helper.ConvertStringToUUID(request.UserUUID)
	if err != nil {
		return response.AccountResponseCreate{}, err
	}

	// check if user already exists
	user, err := a.db.GetUserByUserUUID(ctx, useruuid)
	if err != nil {
		return response.AccountResponseCreate{}, err
	}

	// check if account already exists
	param := db.GetAccountByUserUUIDAndCurrencyParams{
		UserUuid: user.UserUuid,
		Currency: request.Currency,
	}

	checkUserUUIDCurrency, err := a.db.GetAccountByUserUUIDAndCurrency(ctx, param)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return response.AccountResponseCreate{}, err
		}
	}

	if checkUserUUIDCurrency.ID != 0 {
		// return error account already exists
		return response.AccountResponseCreate{}, nil
	}
	// create account
	account, err := a.db.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    request.Owner,
		Currency: request.Currency,
		Balance:  "0",
		UserUuid: useruuid,
	})

	if err != nil {
		return response.AccountResponseCreate{}, err
	}

	result := response.AccountResponseCreate{
		AccountUUID: account.AccountUuid,
		Owner:       account.Owner,
		Currency:    account.Currency,
		Balance:     account.Balance,
		User: response.UserGetSimple{
			UserUUID: user.UserUuid.String(),
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		},
	}

	return result, nil
}

func (a *AccountService) GetAccountByUUID(ctx context.Context, uuid uuid.UUID) (response.AccountResponseGet, error) {
	account, err := a.db.GetAccountByUUID(ctx, uuid)

	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return response.AccountResponseGet{}, err
		}

	}

	user, err := a.db.GetUserByUserUUID(ctx, account.UserUuid)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return response.AccountResponseGet{}, err
		}
	}

	result := response.AccountResponseGet{
		AccountUUID: account.AccountUuid,
		Owner:       account.Owner,
		Currency:    account.Currency,
		Balance:     account.Balance,
		CreatedAt:   account.CreatedAt,
		Status:      int32(account.Status),
		User: response.UserGetSimple{
			UserUUID: user.UserUuid.String(),
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		},
	}

	return result, nil
}

func (a *AccountService) ListAccount(ctx context.Context, param db.ListAccountsParams) ([]response.AccountResponseGet, int64, error) {
	accounts, err := a.db.ListAccounts(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	// count total Data
	countTotal, err := a.db.CountAccounts(ctx)
	if err != nil {
		return nil, 0, err
	}

	var result []response.AccountResponseGet
	for _, account := range accounts {
		// if i query to get user by useruuid here is n+1 problem
		result = append(result, response.AccountResponseGet{
			AccountUUID: account.AccountUuid,
			Owner:       account.Owner,
			Currency:    account.Currency,
			Balance:     account.Balance,
			CreatedAt:   account.CreatedAt,
			Status:      int32(account.Status),
			User: response.UserGetSimple{
				UserUUID: account.UserUuid.String(),
				Username: account.Username.String,
				FullName: account.FullName.String,
				Email:    account.Email.String,
			},
		})
	}

	return result, countTotal, nil
}

func (a *AccountService) UpdateAccount(ctx context.Context, arg db.UpdateProfileAccountParams) (response.AccountResponseGet, error) {
	account, err := a.db.UpdateProfileAccount(ctx, arg)
	if err != nil {
		return response.AccountResponseGet{}, err
	}

	result := response.AccountResponseGet{
		AccountUUID: account.AccountUuid,
		Owner:       account.Owner,
		Currency:    account.Currency,
		Balance:     account.Balance,
		CreatedAt:   account.CreatedAt,
		Status:      int32(account.Status),
		User: response.UserGetSimple{
			UserUUID: account.UserUuid.String(),
			Username: account.Username.String,
			FullName: account.FullName.String,
			Email:    account.Email.String,
		},
	}

	return result, nil
}
