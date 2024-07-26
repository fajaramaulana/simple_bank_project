package response

type UserGetSimple struct {
	UserUUID string `json:"user_uuid"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type UserResponseCreate struct {
	UserUUID string                `json:"user_uuid"`
	FullName string                `json:"full_name"`
	Email    string                `json:"email"`
	Username string                `json:"username"`
	Account  AccountResponseSimple `json:"account"`
}
