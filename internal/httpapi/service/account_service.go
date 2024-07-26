package service

import (
	"context"
	"errors"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/response"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
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

// CreateAccount creates a new account for the authenticated user.
// It takes the context, a CreateAccountRequest object containing the account details,
// and the authentication payload as input parameters.
// It returns an AccountResponseCreate object containing the created account details,
// or an error if the account creation fails.
func (a *AccountService) CreateAccount(ctx context.Context, request *request.CreateAccountRequest, authPayload *token.Payload) (response.AccountResponseCreate, error) {
	useruuid := authPayload.UserUUID

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
		Owner:    user.FullName,
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

func (a *AccountService) GetAccountByUUID(ctx context.Context, uuid uuid.UUID, authPayload *token.Payload) (response.AccountResponseGet, error) {
	account, err := a.db.GetAccountByUUID(ctx, uuid)
	if err != nil {
		// if err.Error() != "sql: no rows in result set" {
		return response.AccountResponseGet{}, err
		// }

	}

	if account.UserUuid != authPayload.UserUUID {
		return response.AccountResponseGet{}, errors.New("unauthorized")
	}

	user, err := a.db.GetUserByUserUUID(ctx, account.UserUuid)
	if err != nil {
		// if err.Error() != "sql: no rows in result set" {
		return response.AccountResponseGet{}, err
		// }
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

func (a *AccountService) ListAccount(ctx context.Context, param db.ListAccountsParams, authPayload *token.Payload) ([]response.AccountResponseGet, int64, error) {
	// if role is admin or superadmin return all account
	if authPayload.Role == "admin" || authPayload.Role == "superadmin" {
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
	} else {

		param := db.ListAccountsByUserUUIDParams{
			UserUuid: authPayload.UserUUID,
			Limit:    param.Limit,
			Offset:   param.Offset,
		}
		accounts, err := a.db.ListAccountsByUserUUID(ctx, param)
		if err != nil {
			return nil, 0, err
		}

		// count total Data
		countTotal, err := a.db.CountAccountsByUserUUID(ctx, authPayload.UserUUID)
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

}

func (a *AccountService) UpdateAccount(ctx context.Context, arg db.UpdateProfileAccountParams, authPayload *token.Payload) (response.AccountResponseGet, error) {
	// check if role is not admin or superadmin
	if authPayload.Role != "admin" && authPayload.Role != "superadmin" {
		// check if account is not owned by the user
		account, err := a.db.GetAccountByUUID(ctx, arg.AccountUuid)
		if err != nil {
			return response.AccountResponseGet{}, err
		}

		if account.UserUuid != authPayload.UserUUID {
			return response.AccountResponseGet{}, errors.New("unauthorized")
		}
	}

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
