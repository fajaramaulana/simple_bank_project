package db

import (
	"context"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store represents the interface for interacting with the database.
type Store interface {
	CreateUserWithAccountTx(ctx context.Context, arg request.CreateUserRequest) (CreateUserWithAccountResult, error)
	TransferTx(ctx context.Context, param TransferTxParam) (TransferTxResult, error)
	Querier
}

// SQLStore provides all functions to execute SQL queries and transactions.
type SQLStore struct {
	conPool *pgxpool.Pool
	*Queries
}

type TransferTxParam struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Type          string `json:"type"`
}

type TransferTxResult struct {
	Transaction Transaction               `json:"transaction"`
	FromAccount SubtractAccountBalanceRow `json:"from_account"`
	ToAccount   AddAccountBalanceRow      `json:"to_account"`
	FromEntry   Entry                     `json:"from_entry"`
	ToEntry     Entry                     `json:"to_entry"`
}

type CreateUserWithAccountResult struct {
	User    response.UserGetSimple         `json:"user"`
	Account response.AccountResponseSimple `json:"account"`
}

// NewStore creates a new instance of the Store interface.
// It takes a *sql.DB as a parameter and returns a Store.
func NewStore(conPool *pgxpool.Pool) Store {
	return &SQLStore{
		conPool: conPool,
		Queries: New(conPool),
	}
}
