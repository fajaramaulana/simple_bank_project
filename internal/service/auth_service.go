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

func (a *AuthService) Login(ctx context.Context, username, password, userAgent, ClientIP string) (response.AuthLoginResponse, error) {
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
	accessTokenDuration, err := time.ParseDuration(a.configToken["access_token_duration"])
	if err != nil {
		return response.AuthLoginResponse{}, err
	}

	role := detailLogin.Role

	accessToken, _, err := maker.CreateToken(detailLogin.UserUuid.String(), accessTokenDuration, role)
	if err != nil {
		return response.AuthLoginResponse{}, err
	}

	refreshTokenDuration, err := time.ParseDuration(a.configToken["refresh_token_duration"])
	if err != nil {
		return response.AuthLoginResponse{}, err
	}

	refreshToken, payloadRefresh, err := maker.CreateToken(detailLogin.UserUuid.String(), refreshTokenDuration, role)
	if err != nil {
		return response.AuthLoginResponse{}, err
	}

	// save refresh token to db
	session, err := a.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           payloadRefresh.ID,
		UserUuid:     detailLogin.UserUuid,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     ClientIP,
		IsBlocked:    false,
		ExpiresAt:    payloadRefresh.ExpiredAt,
	})

	if err != nil {
		return response.AuthLoginResponse{}, err
	}

	return response.AuthLoginResponse{
		SessionId:    session.ID.String(),
		AcessToken:   accessToken,
		RefreshToken: refreshToken,
		User: response.UserGetSimple{
			UserUUID: detailLogin.UserUuid.String(),
			Username: detailLogin.Username,
			FullName: detailLogin.FullName,
			Email:    detailLogin.Email,
		},
	}, nil
}
