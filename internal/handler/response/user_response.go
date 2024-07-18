package response

type UserGetSimple struct {
	UserUUID string `json:"user_uuid"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
