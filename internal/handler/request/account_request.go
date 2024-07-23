package request

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type GetAccountRequest struct {
	UUIDAcc string `uri:"uuid" binding:"required"`
}

type ListAccountRequest struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Limit int32 `form:"limit" binding:"required,min=5,max=10"`
}

type UpdateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR IDR"`
	Status   int32  `json:"status" binding:"required"`
}
