package token_test

import (
	"testing"
	"time"

	"github.com/fajaramaulana/simple_bank_project/internal/handler/token"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	newRandomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	// parse uuid to string
	uuidUser := newRandomUUID.String()
	// role := util.RandomRole()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(uuidUser, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, uuidUser, payload.UserUUID.String())
	// require.Equal(t, role, payload.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	newRandomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	// parse uuid to string
	uuidUser := newRandomUUID.String()

	tokenString, payload, err := maker.CreateToken(uuidUser, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrExpiredToken.Error())
	require.Nil(t, payload)
}
func TestInvalidTokenPaseto(t *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	_, err = maker.VerifyToken("invalid token")
	require.EqualError(t, err, token.ErrInvalidToken.Error())
}

func TestInvalidPayloadPaseto(t *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	_, _, err = maker.CreateToken("test", time.Minute)
	require.Error(t, err)
}

func TestInvalidSecretKeyPaseto(t *testing.T) {
	_, err := token.NewPasetoMaker("short")
	require.EqualError(t, err, "invalid key, must be at least 32 characters")
}
