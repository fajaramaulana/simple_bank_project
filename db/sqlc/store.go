package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type Store struct {
	*Queries
	db *sql.DB
}

// NewStore create a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

type TransferTxParam struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Type          string `json:"type"`
}

type TransferTxResult struct {
	Transaction Transaction `json:"transaction"`
	FromAccount Account     `json:"from_account"`
	ToAccount   Account     `json:"to_account"`
	FromEntry   Entry       `json:"from_entry"`
	ToEntry     Entry       `json:"to_entry"`
}

// TransferTx perform money transfer from one account to another one account
func (store *Store) TransferTx(ctx context.Context, param TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			FromAccountID: param.FromAccountID,
			ToAccountID:   param.ToAccountID,
			Amount:        strconv.FormatInt(param.Amount, 10),
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.FromAccountID,
			Amount:    strconv.FormatInt(param.Amount, 10),
			TypeTrans: "debit",
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.ToAccountID,
			Amount:    strconv.FormatInt(param.Amount, 10),
			TypeTrans: "credit",
		})
		if err != nil {
			return err
		}

		// TODO update account balance
		account1, err := q.GetAccountForUpdate(ctx, param.FromAccountID)

		if err != nil {
			return err
		}

		account1BalanceConvFloat64, err := strconv.ParseFloat(account1.Balance, 64)
		if err != nil {
			return err
		}

		updateBalanceFromAccount := account1BalanceConvFloat64 - float64(param.Amount)

		updateBalanceFromAccountNew := strconv.FormatFloat(updateBalanceFromAccount, 'f', 6, 64)

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      param.FromAccountID,
			Balance: updateBalanceFromAccountNew,
		})

		if err != nil {
			return err
		}

		account2, err := q.GetAccountForUpdate(ctx, param.ToAccountID)
		if err != nil {
			return err
		}

		account2BalanceConvFloat64, err := strconv.ParseFloat(account2.Balance, 64)
		if err != nil {
			return err
		}

		updateBalanceToAccount := account2BalanceConvFloat64 + float64(param.Amount)

		updateBalanceToAccountNew := strconv.FormatFloat(updateBalanceToAccount, 'f', 6, 64)

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      param.ToAccountID,
			Balance: updateBalanceToAccountNew,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
