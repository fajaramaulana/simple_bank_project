package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/response"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/token"
	"github.com/fajaramaulana/simple_bank_project/util"
)

var (
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrorInvalidPassword = fmt.Errorf("invalid password")
)

type AuthService struct {
	db          db.Store
	configToken map[string]string
}

func NewAuthService(db db.Store, configToken map[string]string) *AuthService {
	return &AuthService{
		db:          db,
		configToken: configToken,
	}
}

func (a *AuthService) Login(ctx context.Context, username, password string) (response.AuthLoginResponse, error) {
	detailLogin, err := a.db.GetDetailLoginByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.AuthLoginResponse{}, ErrUserNotFound
		}
		return response.AuthLoginResponse{}, err
	}

	// check password
	err = util.CheckPasswordBcrypt(password, detailLogin.HashedPassword)
	if err != nil {
		return response.AuthLoginResponse{}, ErrorInvalidPassword
	}

	// generate token with paseto
	maker, err := token.NewPasetoMaker(a.configToken["token_secret"])
	if err != nil {
		return response.AuthLoginResponse{}, err
	}
	tokenDuration, err := time.ParseDuration(a.configToken["access_token_duration"])
	if err != nil {
		return response.AuthLoginResponse{}, err
	}
	duration := tokenDuration

	role := detailLogin.Role

	token, _, err := maker.CreateToken(detailLogin.UserUuid.String(), duration, role)
	if err != nil {
		return response.AuthLoginResponse{}, err
	}

	return response.AuthLoginResponse{
		AcessToken: token,
		User: response.UserGetSimple{
			UserUUID: detailLogin.UserUuid.String(),
			Username: detailLogin.Username,
			FullName: detailLogin.FullName,
			Email:    detailLogin.Email,
		},
	}, nil
}
