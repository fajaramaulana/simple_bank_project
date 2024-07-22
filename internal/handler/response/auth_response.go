package response

type AuthLoginResponse struct {
	AcessToken string        `json:"access_token"`
	User       UserGetSimple `json:"user"`
}
