package request

type AuthLoginRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=6"`
	Password string `json:"password" binding:"required,min=8"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
