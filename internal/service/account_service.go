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
	db *db.Store
}

func NewAccountService(db *db.Store) *AccountService {
	return &AccountService{
		db: db,
	}
}

func (a *AccountService) CreateAccount(ctx context.Context, request *request.CreateAccountRequest) (response.AccountResponseCreate, error) {
	// check if account already exists
	checkEmail, err := a.db.GetAccountByEmail(ctx, request.Email)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return response.AccountResponseCreate{}, err
		}
	}

	if checkEmail.ID != 0 {
		// return error account already exists
		return response.AccountResponseCreate{}, nil
	}
	// create account

	// generate password bcrypt
	password, err := helper.GeneratePasswordBcrypt(request.Password)
	if err != nil {
		return response.AccountResponseCreate{}, err
	}

	account, err := a.db.CreateAccount(ctx, db.CreateAccountParams{
		Owner:        request.Owner,
		Email:        request.Email,
		Password:     password,
		Currency:     "USD",
		Balance:      "0",
		RefreshToken: "refresh_token",
	})

	if err != nil {
		return response.AccountResponseCreate{}, err
	}

	result := response.AccountResponseCreate{
		AccountUUID: account.AccountUuid,
		Owner:       account.Owner,
		Email:       account.Email,
		Currency:    account.Currency,
		Balance:     account.Balance,
	}

	return result, nil
}

func (a *AccountService) GetAccountByUUID(ctx context.Context, uuid uuid.UUID) (response.AccountResponseGet, error) {
	account, err := a.db.GetAccountByUUID(ctx, uuid)

	if err != nil {
		return response.AccountResponseGet{}, err
	}

	result := response.AccountResponseGet{
		AccountUUID: account.AccountUuid,
		Owner:       account.Owner,
		Email:       account.Email,
		Currency:    account.Currency,
		Balance:     account.Balance,
		CreatedAt:   account.CreatedAt.String(),
		Status:      account.Status,
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
		result = append(result, response.AccountResponseGet{
			AccountUUID: account.AccountUuid,
			Owner:       account.Owner,
			Email:       account.Email,
			Currency:    account.Currency,
			Balance:     account.Balance,
			CreatedAt:   account.CreatedAt.String(),
			Status:      account.Status,
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
		Email:       account.Email,
		Currency:    account.Currency,
		Balance:     account.Balance,
		CreatedAt:   account.CreatedAt.String(),
		Status:      account.Status,
	}

	return result, nil
}
