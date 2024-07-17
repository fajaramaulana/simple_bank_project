package request

type CreateTransferRequest struct {
	FromAccountUUID string `json:"from_account_uuid" binding:"required"`
	ToAccountUUID   string `json:"to_account_uuid" binding:"required"`
	Amount          int64  `json:"amount" binding:"required,numeric,min=1"`
	Currency        string `json:"currency" binding:"required,currency"`
}
