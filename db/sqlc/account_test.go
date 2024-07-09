package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/stretchr/testify/require"
)

func generateAccount(t *testing.T) Account {

	r := util.NewRandomMoneyGenerator()
	input := CreateAccountParams{
		Owner:    util.RandomName(),
		Balance:  util.RandomMoney(r, 10.00, 99999999.00),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), input)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, input.Owner, account.Owner, "input and return owner should be same")
	require.Equal(t, input.Balance, account.Balance, "input and return balance should be same")
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

	getFromAccount, err := testQueries.GetAccount(context.Background(), createRandAccount.ID)

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

	updatedData, err := testQueries.UpdateAccount(context.Background(), params)
	require.NoError(t, err)

	require.NotEmpty(t, updatedData)
	require.Equal(t, params.ID, updatedData.ID)
	require.Equal(t, params.Balance, updatedData.Balance)
	require.NotEmpty(t, updatedData.UpdatedAt)
}

func TestSoftDeleteAccount(t *testing.T) {
	createRandAccount := generateAccount(t)

	err := testQueries.SoftDeleteAccount(context.Background(), createRandAccount.ID)
	require.NoError(t, err)
	getAccount, err := testQueries.GetAccount(context.Background(), createRandAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getAccount)

}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = generateAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
