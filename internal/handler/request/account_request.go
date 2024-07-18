package request

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
	UserUUID string `json:"user_uuid" binding:"required"`
	Password string `json:"password" binding:"required"`
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
