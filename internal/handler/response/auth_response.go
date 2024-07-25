package response

type AuthLoginResponse struct {
	SessionId    string        `json:"session_id"`
	AcessToken   string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         UserGetSimple `json:"user"`
}
