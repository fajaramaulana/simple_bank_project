package request

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=6"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}
