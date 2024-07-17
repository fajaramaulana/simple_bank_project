package response

type SuccessTransactionResponse struct {
	TransactionUUID string `json:"transaction_uuid"`
	FromAccountUUID string `json:"from_account_uuid"`
	ToAccountUUID   string `json:"to_account_uuid"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	LastedBalance   string `json:"lasted_balance"`
	Type            string `json:"type"`
}
