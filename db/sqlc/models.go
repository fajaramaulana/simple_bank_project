// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID          int64        `json:"id"`
	Owner       string       `json:"owner"`
	Currency    string       `json:"currency"`
	Balance     string       `json:"balance"`
	CreatedAt   time.Time    `json:"created_at"`
	AccountUuid uuid.UUID    `json:"account_uuid"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	// can be negative or positive
	Amount      string       `json:"amount"`
	CreatedAt   time.Time    `json:"created_at"`
	EntriesUuid uuid.UUID    `json:"entries_uuid"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}

type Transaction struct {
	ID              int64        `json:"id"`
	FromAccountID   int64        `json:"from_account_id"`
	ToAccountID     int64        `json:"to_account_id"`
	Amount          string       `json:"amount"`
	CreatedAt       time.Time    `json:"created_at"`
	TransactionUuid uuid.UUID    `json:"transaction_uuid"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at"`
}
