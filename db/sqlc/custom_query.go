package db

import (
	"context"
	"fmt"
	"math/big"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/response"
	"github.com/jackc/pgx/v5/pgtype"
)

// execTx executes the given function within a transaction.
// It begins a new transaction, calls the provided function with a `Queries` instance,
// and commits the transaction if the function returns nil.
// If the function returns an error, it rolls back the transaction and returns the error.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// tx, err := store.db.BeginTx(ctx, nil)
	tx, err := store.conPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

// TransferTx perform money transfer from one account to another one account
// TransferTx performs a transaction that transfers funds from one account to another.
// It creates a transaction record, debit entry for the sender account, credit entry for the receiver account,
// and updates the account balances accordingly.
// The function takes a context and a TransferTxParam as input and returns a TransferTxResult and an error.
func (store *SQLStore) TransferTx(ctx context.Context, param TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		arg := CreateTransactionParams{
			FromAccountID: param.FromAccountID,
			ToAccountID:   param.ToAccountID,
			Amount:        pgtype.Numeric{Int: big.NewInt(param.Amount), Exp: 0, Valid: true},
		}
		fmt.Printf("%# v\n", arg)
		result.Transaction, err = q.CreateTransaction(ctx, arg)

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.FromAccountID,
			Amount:    pgtype.Numeric{Int: big.NewInt(param.Amount), Exp: 0, Valid: true},
			TypeTrans: "debit",
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.ToAccountID,
			Amount:    pgtype.Numeric{Int: big.NewInt(param.Amount), Exp: 0, Valid: true},
			TypeTrans: "credit",
		})
		if err != nil {
			return err
		}

		// TODO update account balance
		result.FromAccount, result.ToAccount, err = addBalance(ctx, q, param.FromAccountID, param.Amount, param.ToAccountID, param.Amount)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (store *SQLStore) CreateUserWithAccountTx(ctx context.Context, arg request.CreateUserRequest) (CreateUserWithAccountResult, error) {
	var result CreateUserWithAccountResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		user, err := q.CreateUser(ctx, CreateUserParams{
			Username:       arg.Username,
			FullName:       arg.FullName,
			HashedPassword: arg.Password,
			Email:          arg.Email,
		})

		if err != nil {
			return err
		}

		account, err := q.CreateAccount(ctx, CreateAccountParams{
			Owner:    user.FullName,
			Balance:  pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true},
			Currency: arg.Currency,
			UserUuid: user.UserUuid,
		})

		if err != nil {
			return err
		}

		result.Account = response.AccountResponseSimple{
			AccountUUID: account.AccountUuid,
			Owner:       account.Owner,
			Currency:    account.Currency,
			Balance:     account.Balance.Int.String(),
		}

		result.User = response.UserGetSimple{
			UserUUID: user.UserUuid.String(),
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		}

		return nil
	})

	return result, err
}

// addBalance subtracts the specified amount from the account with accountID1 and adds the specified amount to the account with accountID2.
// If accountID1 is less than accountID2, the subtraction operation is performed first, followed by the addition operation.
// If accountID1 is greater than or equal to accountID2, the addition operation is performed first, followed by the subtraction operation.
// The function returns the updated account balances for both accounts and any error that occurred during the operations.
func addBalance(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (fromAccount SubtractAccountBalanceRow, toAccount AddAccountBalanceRow, err error) {
	if accountID1 < accountID2 {
		fromAccount, err := q.SubtractAccountBalance(ctx, SubtractAccountBalanceParams{
			Amount: pgtype.Numeric{Int: big.NewInt(amount1), Exp: 0, Valid: true},
			ID:     accountID1,
		})

		if err != nil {
			return fromAccount, toAccount, err
		}

		toAccount, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: pgtype.Numeric{Int: big.NewInt(amount2), Exp: 0, Valid: true},
			ID:     accountID2,
		})

		if err != nil {
			return fromAccount, toAccount, err
		}

		return fromAccount, toAccount, nil
	} else {
		toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: pgtype.Numeric{Int: big.NewInt(amount2), Exp: 0, Valid: true},
			ID:     accountID2,
		})

		if err != nil {
			return fromAccount, toAccount, err
		}

		fromAccount, err := q.SubtractAccountBalance(ctx, SubtractAccountBalanceParams{
			Amount: pgtype.Numeric{Int: big.NewInt(amount1), Exp: 0, Valid: true},
			ID:     accountID1,
		})

		if err != nil {
			return fromAccount, toAccount, err
		}

		return fromAccount, toAccount, nil
	}
}
