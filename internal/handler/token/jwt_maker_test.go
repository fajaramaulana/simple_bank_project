package token_test

import (
	"testing"
	"time"

	"github.com/fajaramaulana/simple_bank_project/internal/handler/token"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := token.NewJWTMaker(util.RandomString(32))
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

func TestExpiredJWTToken(t *testing.T) {
	maker, err := token.NewJWTMaker(util.RandomString(32))
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

func TestInvalidJWTTokenAlgNone(t *testing.T) {

	newRandomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	// parse uuid to string
	uuidUser := newRandomUUID.String()

	payload, err := token.NewPayload(uuidUser, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	tokenString, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestInvalidSecretToken(t *testing.T) {
	maker, err := token.NewJWTMaker(util.RandomString(30))
	require.Error(t, err)
	require.Nil(t, maker)
}

func TestInvalidPayload(t *testing.T) {
	_, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	_, err = token.NewPayload("wrong-uuid", time.Minute)
	require.Error(t, err)
	require.EqualError(t, err, "invalid UUID length: 10")
}
