package response

import (
	"time"

	"github.com/google/uuid"
)

type AccountResponseCreate struct {
	AccountUUID uuid.UUID     `json:"account_uuid"`
	Owner       string        `json:"owner"`
	Currency    string        `json:"currency"`
	Balance     string        `json:"balance"`
	User        UserGetSimple `json:"user"`
}

type AccountResponseGet struct {
	AccountUUID uuid.UUID     `json:"account_uuid"`
	Owner       string        `json:"owner"`
	Currency    string        `json:"currency"`
	Balance     string        `json:"balance"`
	CreatedAt   time.Time     `json:"created_at"`
	Status      int32         `json:"status"`
	User        UserGetSimple `json:"user"`
}
