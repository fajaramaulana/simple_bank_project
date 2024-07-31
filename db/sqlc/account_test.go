package db

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func generateAccount(t *testing.T) CreateAccountRow {
	password, err := util.MakePasswordBcrypt("P4ssw0rd!")
	require.NoError(t, err)
	inputUser := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: password,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), inputUser)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, inputUser.Username, user.Username, "input and return username should be same")
	require.Equal(t, inputUser.Email, user.Email, "input and return email should be same")
	require.NotEmpty(t, user.UserUuid.String())
	require.NotEmpty(t, user.Role)

	r := util.NewRandomMoneyGenerator()
	input := CreateAccountParams{
		Owner:    user.FullName,
		Balance:  pgtype.Numeric{Int: big.NewInt(util.RandomMoneyInt(r, 10, 99999999)), Exp: 0, Valid: true},
		Currency: util.RandomCurrency(),
		UserUuid: user.UserUuid,
	}

	account, err := testStore.CreateAccount(context.Background(), input)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	inputBalanceBigInt := helper.NumericToBigInt(input.Balance)
	accountBalanceBigInt := helper.NumericToBigInt(account.Balance)

	// Compare the balances
	require.Equal(t, input.Owner, account.Owner, "input and return owner should be same")
	require.Equal(t, inputBalanceBigInt.String(), accountBalanceBigInt.String(), "input and return balance should be same")
	require.Equal(t, input.Currency, account.Currency, "input and return currency should be same")
	require.NotEmpty(t, account.AccountUuid.String())
	require.NotNil(t, account.CreatedAt)
	require.NotZero(t, account.ID)

	return account
}

func TestCreateAccount(t *testing.T) {
	generateAccount(t)
}

func TestGetAccount(t *testing.T) {
	createRandAccount := generateAccount(t)

	getFromAccount, err := testStore.GetAccount(context.Background(), createRandAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getFromAccount)
	require.Equal(t, createRandAccount.ID, getFromAccount.ID, "create and get should be")
	require.Equal(t, createRandAccount.AccountUuid.String(), getFromAccount.AccountUuid.String(), "create and get should be")
	require.Equal(t, createRandAccount.Balance, getFromAccount.Balance, "create and get should be")
	require.Equal(t, createRandAccount.Currency, getFromAccount.Currency, "create and get should be")
	require.WithinDuration(t, createRandAccount.CreatedAt, getFromAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	createRandAccount := generateAccount(t)

	params := UpdateAccountParams{
		ID:      createRandAccount.ID,
		Balance: createRandAccount.Balance,
	}

	updatedData, err := testStore.UpdateAccount(context.Background(), params)
	require.NoError(t, err)

	require.NotEmpty(t, updatedData)
	require.Equal(t, params.ID, updatedData.AccountID)
	require.Equal(t, params.Balance, updatedData.Balance)
	require.NotEmpty(t, updatedData.UpdatedAt)
}

func TestSoftDeleteAccount(t *testing.T) {
	createRandAccount := generateAccount(t)

	err := testStore.SoftDeleteAccount(context.Background(), createRandAccount.ID)
	require.NoError(t, err)
	getAccount, err := testStore.GetAccount(context.Background(), createRandAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, getAccount)

}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		generateAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestListAccountsNil(t *testing.T) {
	var lastAccount CreateAccountRow
	for i := 0; i < 10; i++ {
		lastAccount = generateAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 444,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, accounts)

	for _, account := range accounts {
		require.Empty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestGetAccountForUpdate(t *testing.T) {
	createRandAccount := generateAccount(t)

	getFromAccount, err := testStore.GetAccountForUpdate(context.Background(), createRandAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getFromAccount)
	require.Equal(t, createRandAccount.ID, getFromAccount.ID, "create and get should be")
	require.Equal(t, createRandAccount.AccountUuid.String(), getFromAccount.AccountUuid.String(), "create and get should be")
}

func TestListAccountsError(t *testing.T) {
	for i := 0; i < 10; i++ {
		generateAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  -5,
		Offset: 0,
	}

	_, err := testStore.ListAccounts(context.Background(), arg)
	require.Error(t, err)
}

func TestGetAccountByUUID(t *testing.T) {
	createRandAccount := generateAccount(t)

	getFromAccount, err := testStore.GetAccountByUUID(context.Background(), createRandAccount.AccountUuid)

	require.NoError(t, err)
	require.NotEmpty(t, getFromAccount)
	require.Equal(t, createRandAccount.ID, getFromAccount.ID, "create and get should be")
	require.Equal(t, createRandAccount.AccountUuid.String(), getFromAccount.AccountUuid.String(), "create and get should be")
	require.Equal(t, createRandAccount.Balance, getFromAccount.Balance, "create and get should be")
	require.Equal(t, createRandAccount.Currency, getFromAccount.Currency, "create and get should be")
	require.WithinDuration(t, createRandAccount.CreatedAt, getFromAccount.CreatedAt, time.Second)
}
