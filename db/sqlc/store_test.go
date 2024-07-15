package db

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := generateAccount(t)
	account2 := generateAccount(t)

	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// make random ["credit", "debit"]
	transactionType := getRandomOption([]string{"credit", "debit"})

	// run n concurent transfer transaction
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParam{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Type:          transactionType,
			})

			errs <- err
			results <- result

		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results

		transaction := result.Transaction
		// fmt.Printf("%# v\n", pretty.Formatter(transaction))

		convertAmontFloat64, err := strconv.ParseFloat(transaction.Amount, 64)

		require.NoError(t, err)

		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, account1.ID, transaction.FromAccountID)
		require.Equal(t, account2.ID, transaction.ToAccountID)
		require.Equal(t, float64(amount), convertAmontFloat64)
		require.NotZero(t, transaction.ID)
		require.NotZero(t, transaction.CreatedAt)
		require.NotZero(t, transaction.ToAccountID)

		_, err = store.GetTransaction(context.Background(), transaction.ID)
		require.NoError(t, err)

		// check entries
		convertAmountFromEntryFloat64, err := strconv.ParseFloat(result.FromEntry.Amount, 64)
		require.NoError(t, err)
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, float64(amount), convertAmountFromEntryFloat64)
		require.Equal(t, "debit", fromEntry.TypeTrans)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		convertAmountToEntryFloat64, err := strconv.ParseFloat(result.ToEntry.Amount, 64)

		require.NoError(t, err)
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, float64(amount), convertAmountToEntryFloat64)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balance
		account1Balance, err := strconv.ParseFloat(account1.Balance, 64)
		require.NoError(t, err)
		fromAccountBalance, err := strconv.ParseFloat(fromAccount.Balance, 64)
		require.NoError(t, err)
		toAccountBalance, err := strconv.ParseFloat(toAccount.Balance, 64)
		require.NoError(t, err)

		account2Balance, err := strconv.ParseFloat(account2.Balance, 64)
		require.NoError(t, err)
		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1Balance - fromAccountBalance
		diff2 := toAccountBalance - account2Balance
		fmt.Printf("%# v\n", diff1)
		fmt.Printf("%# v\n", diff2)

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, int64(diff1)%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(int64(diff1) / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	accountOneBalance, err := strconv.ParseFloat(account1.Balance, 64)
	require.NoError(t, err)
	accountTwoBalance, err := strconv.ParseFloat(account2.Balance, 64)
	require.NoError(t, err)

	// convert updatedAccount1.Balance string to float64
	updatedAccount1Balance, err := strconv.ParseFloat(updatedAccount1.Balance, 64)
	require.NoError(t, err)

	// convert updatedAccount2.Balance string to float64
	updatedAccount2Balance, err := strconv.ParseFloat(updatedAccount2.Balance, 64)
	require.NoError(t, err)

	require.Equal(t, accountOneBalance-float64(n)*float64(amount), updatedAccount1Balance)
	require.Equal(t, accountTwoBalance+float64(n)*float64(amount), updatedAccount2Balance)
}

func TestTransTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := generateAccount(t)
	account2 := generateAccount(t)

	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// make random ["credit", "debit"]
	transactionType := getRandomOption([]string{"credit", "debit"})

	// run n concurent transfer transaction
	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {

		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParam{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
				Type:          transactionType,
			})

			errs <- err

		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	accountOneBalance, err := strconv.ParseFloat(account1.Balance, 64)
	require.NoError(t, err)
	accountTwoBalance, err := strconv.ParseFloat(account2.Balance, 64)
	require.NoError(t, err)

	// convert updatedAccount1.Balance string to float64
	updatedAccount1Balance, err := strconv.ParseFloat(updatedAccount1.Balance, 64)
	require.NoError(t, err)

	// convert updatedAccount2.Balance string to float64
	updatedAccount2Balance, err := strconv.ParseFloat(updatedAccount2.Balance, 64)
	require.NoError(t, err)

	require.Equal(t, accountOneBalance, updatedAccount1Balance)
	require.Equal(t, accountTwoBalance, updatedAccount2Balance)
}

func getRandomOption(options []string) string {
	n := len(options)
	idx := rand.Intn(n)
	return options[idx]
}
