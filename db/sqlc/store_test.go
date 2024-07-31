package db

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	helpergrpc "github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/stretchr/testify/require"
)

func TestTransTx(t *testing.T) {
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
			result, err := testStore.TransferTx(context.Background(), TransferTxParam{
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

		getData, err := transaction.Amount.Float64Value()
		require.NoError(t, err)
		convertAmount := helpergrpc.NumerictoFloat64(transaction.Amount)
		fmt.Printf("%# v\n", getData)

		require.NoError(t, err)

		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, account1.ID, transaction.FromAccountID)
		require.Equal(t, account2.ID, transaction.ToAccountID)
		require.Equal(t, float64(amount), convertAmount)
		require.NotZero(t, transaction.ID)
		require.NotZero(t, transaction.CreatedAt)
		require.NotZero(t, transaction.ToAccountID)

		_, err = testStore.GetTransaction(context.Background(), transaction.ID)
		require.NoError(t, err)

		// check entries
		convertAmountFromEntry := helpergrpc.NumerictoFloat64(result.FromEntry.Amount)
		require.NoError(t, err)
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, float64(amount), convertAmountFromEntry)
		require.Equal(t, "debit", fromEntry.TypeTrans)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		convertAmountToEntry := helpergrpc.NumerictoFloat64(result.ToEntry.Amount)

		require.NoError(t, err)
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, float64(amount), convertAmountToEntry)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balance
		account1Balance := helpergrpc.NumericToBigInt(account1.Balance)
		fromAccountBalance := helpergrpc.NumericToBigInt(fromAccount.Balance)
		toAccountBalance := helpergrpc.NumericToBigInt(toAccount.Balance)
		account2Balance := helpergrpc.NumericToBigInt(account2.Balance)
		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := new(big.Int).Sub(account1Balance, fromAccountBalance)
		diff2 := new(big.Int).Sub(toAccountBalance, account2Balance)
		fmt.Printf("%# v\n", diff1)
		fmt.Printf("%# v\n", diff2)

		require.Equal(t, diff1, diff2)
		require.True(t, diff1.Cmp(big.NewInt(0)) > 0)
		require.True(t, diff1.Int64()%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1.Int64() / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	accountOneBalance, err := strconv.ParseFloat(account1.Balance.Int.String(), 64)
	require.NoError(t, err)
	accountTwoBalance, err := strconv.ParseFloat(account2.Balance.Int.String(), 64)
	require.NoError(t, err)

	// convert updatedAccount1.Balance string to float64
	updatedAccount1Balance := helpergrpc.NumerictoFloat64(updatedAccount1.Balance)

	// convert updatedAccount2.Balance string to float64
	updatedAccount2Balance := helpergrpc.NumerictoFloat64(updatedAccount2.Balance)

	require.Equal(t, accountOneBalance-float64(n)*float64(amount), updatedAccount1Balance)
	require.Equal(t, accountTwoBalance+float64(n)*float64(amount), updatedAccount2Balance)
}

func TestCreateUserWithAccountTx(t *testing.T) {
	hashPass, err := util.MakePasswordBcrypt("P4ssw0rd!")

	require.NoError(t, err)
	arg := request.CreateUserRequest{
		Username: util.RandomUsername(),
		Email:    util.RandomEmail(),
		Password: hashPass,
		FullName: util.RandomUsername(),
		Currency: util.RandomCurrency(),
	}

	res, err := testStore.CreateUserWithAccountTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, res.User.UserUUID)
	require.NotZero(t, res.Account.AccountUUID)

	userUUID, err := helper.ConvertStringToUUID(res.User.UserUUID)
	require.NoError(t, err)

	// check user
	checkUser, err := testStore.GetUserByUserUUID(context.Background(), userUUID)
	require.NoError(t, err)
	require.NotEmpty(t, checkUser)

	// check account
	checkAccount, err := testStore.GetAccountByUUID(context.Background(), res.Account.AccountUUID)
	require.NoError(t, err)
	require.NotEmpty(t, checkAccount)

	require.Equal(t, checkUser.UserUuid, checkAccount.UserUuid)
}

func TestCreateUserWithAccountTxExist(t *testing.T) {
	user := GenerateUser(t)
	arg := request.CreateUserRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: "P4ssw0rd!",
		FullName: user.FullName,
		Currency: util.RandomCurrency(),
	}

	_, err := testStore.CreateUserWithAccountTx(context.Background(), arg)
	require.Error(t, err)
}

func TestTransTxDeadlock(t *testing.T) {
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
			_, err := testStore.TransferTx(context.Background(), TransferTxParam{
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
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	accountOneBalance, err := strconv.ParseFloat(account1.Balance.Int.String(), 64)
	require.NoError(t, err)
	accountTwoBalance, err := strconv.ParseFloat(account2.Balance.Int.String(), 64)
	require.NoError(t, err)

	// convert updatedAccount1.Balance string to float64
	updatedAccount1Balance, err := strconv.ParseFloat(updatedAccount1.Balance.Int.String(), 64)
	require.NoError(t, err)

	// convert updatedAccount2.Balance string to float64
	updatedAccount2Balance, err := strconv.ParseFloat(updatedAccount2.Balance.Int.String(), 64)
	require.NoError(t, err)

	require.Equal(t, accountOneBalance, updatedAccount1Balance)
	require.Equal(t, accountTwoBalance, updatedAccount2Balance)
}

func getRandomOption(options []string) string {
	n := len(options)
	idx := rand.Intn(n)
	return options[idx]
}
