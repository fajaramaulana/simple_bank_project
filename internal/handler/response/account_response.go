package response

import (
	"time"

	"github.com/google/uuid"
)

type AccountResponseCreate struct {
	AccountUUID uuid.UUID `json:"account_uuid"`
	Owner       string    `json:"owner"`
	Email       string    `json:"email"`
	Currency    string    `json:"currency"`
	Balance     string    `json:"balance"`
}

type AccountResponseGet struct {
	AccountUUID uuid.UUID `json:"account_uuid"`
	Owner       string    `json:"owner"`
	Email       string    `json:"email"`
	Currency    string    `json:"currency"`
	Balance     string    `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
	Status      int32     `json:"status"`
}
