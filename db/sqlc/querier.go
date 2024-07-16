// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (AddAccountBalanceRow, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	GetAccount(ctx context.Context, id int64) (GetAccountRow, error)
	GetAccountForUpdate(ctx context.Context, id int64) (GetAccountForUpdateRow, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
	GetTransaction(ctx context.Context, id int64) (Transaction, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]ListAccountsRow, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]Transaction, error)
	SoftDeleteAccount(ctx context.Context, id int64) error
	SubtractAccountBalance(ctx context.Context, arg SubtractAccountBalanceParams) (SubtractAccountBalanceRow, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (UpdateAccountRow, error)
}

var _ Querier = (*Queries)(nil)
