package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func generateSession(t *testing.T) Session {
	password, err := util.MakePasswordBcrypt("P4ssw0rd!")
	require.NoError(t, err)
	inputUser := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: password,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), inputUser)
	require.NoError(t, err)

	uuid := uuid.New()

	// get config token

	tokenMaker, err := token.NewPasetoMaker(os.Getenv("TOKEN_SYMMETRIC_KEY"))
	require.NoError(t, err)

	duration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	require.NoError(t, err)

	// create token
	refreshToken, refreshPayload, err := tokenMaker.CreateToken(user.UserUuid.String(), duration, user.Role)
	require.NoError(t, err)

	input := CreateSessionParams{
		ID:           uuid,
		UserUuid:     user.UserUuid,
		RefreshToken: refreshToken,
		UserAgent:    "user-agent",
		ClientIp:     "1.1.1.1",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}

	session, err := testQueries.CreateSession(context.Background(), input)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, input.ID, session.ID, "input and return id should be same")
	require.Equal(t, input.UserUuid, session.UserUuid, "input and return user_uuid should be same")
	require.Equal(t, input.RefreshToken, session.RefreshToken, "input and return refresh_token should be same")
	require.Equal(t, input.UserAgent, session.UserAgent, "input and return user_agent should be same")
	require.Equal(t, input.ClientIp, session.ClientIp, "input and return client_ip should be same")
	require.Equal(t, input.IsBlocked, session.IsBlocked, "input and return is_blocked should be same")
	require.WithinDuration(t, input.ExpiresAt, session.ExpiresAt.Local(), time.Second, "input and return expires_at should be same")
	require.NotNil(t, session.CreatedAt)

	return session
}

func TestCreateSession(t *testing.T) {
	generateSession(t)
}

func TestGetSession(t *testing.T) {
	session := generateSession(t)

	require.NotEmpty(t, session)

	getFromSession, err := testQueries.GetSession(context.Background(), session.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getFromSession)
	require.Equal(t, session.ID, getFromSession.ID, "create and get should be")
	require.Equal(t, session.UserUuid, getFromSession.UserUuid, "create and get should be")
	require.Equal(t, session.RefreshToken, getFromSession.RefreshToken, "create and get should be")
}
