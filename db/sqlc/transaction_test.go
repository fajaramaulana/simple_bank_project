package db

import (
	"context"
	"math/big"
	"testing"

	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func generateRandomTransaction(t *testing.T, acc1 CreateAccountRow, acc2 CreateAccountRow) Transaction {
	r := util.NewRandomMoneyGenerator()
	input := CreateTransactionParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        pgtype.Numeric{Int: big.NewInt(util.RandomMoneyInt(r, 0.001, 999999999.0)), Valid: true},
	}

	transaction, err := testStore.CreateTransaction(context.Background(), input)

	// if e
	require.NoError(t, err)
	require.NotEmpty(t, transaction)
	require.Equal(t, input.FromAccountID, transaction.FromAccountID)
	require.Equal(t, input.ToAccountID, transaction.ToAccountID)
	require.Equal(t, helper.NumerictoFloat64(input.Amount), helper.NumerictoFloat64(transaction.Amount))
	require.NotEmpty(t, transaction.TransactionUuid.String())
	require.NotNil(t, transaction.CreatedAt)
	require.NotZero(t, transaction.ID)

	return transaction
}

func TestCreateTransaction(t *testing.T) {
	acc1 := generateAccount(t)
	acc2 := generateAccount(t)
	generateRandomTransaction(t, acc1, acc2)
}

func TestGetTransaction(t *testing.T) {

	acc1 := generateAccount(t)
	acc2 := generateAccount(t)
	createRandTrans := generateRandomTransaction(t, acc1, acc2)

	getFromTransaction, err := testStore.GetTransaction(context.Background(), createRandTrans.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getFromTransaction)
	require.Equal(t, createRandTrans.ID, getFromTransaction.ID)
	require.Equal(t, createRandTrans.TransactionUuid.String(), getFromTransaction.TransactionUuid.String())
	require.Equal(t, createRandTrans.Amount, getFromTransaction.Amount)
	require.Equal(t, createRandTrans.CreatedAt, getFromTransaction.CreatedAt)

}

func TestListTransaction(t *testing.T) {
	acc1 := generateAccount(t)
	acc2 := generateAccount(t)

	for i := 0; i < 10; i++ {
		generateRandomTransaction(t, acc1, acc2)
	}

	arg := ListTransactionsParams{
		FromAccountID: acc1.ID,
		Limit:         5,
		Offset:        5,
	}

	transactions, err := testStore.ListTransactions(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transactions)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
		require.Equal(t, acc1.ID, transaction.FromAccountID)
		require.Equal(t, acc2.ID, transaction.ToAccountID)
	}
}
