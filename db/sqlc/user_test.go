package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	generateAccount(t)
}

func TestGetByEmail(t *testing.T) {
	user := GenerateUser(t)

	userByEmail, err := testQueries.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, userByEmail)
	require.Equal(t, user.Email, userByEmail.Email)
	require.Equal(t, user.Username, userByEmail.Username)
	require.Equal(t, user.UserUuid, userByEmail.UserUuid)
}

func TestGetByUsername(t *testing.T) {
	user := GenerateUser(t)

	userByUsername, err := testQueries.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userByUsername)
	require.Equal(t, user.Email, userByUsername.Email)
	require.Equal(t, user.Username, userByUsername.Username)
	require.Equal(t, user.UserUuid, userByUsername.UserUuid)
}

func TestUpdateUser(t *testing.T) {
	user := GenerateUser(t)

	arg := UpdateUserParams{
		UserUuid:          user.UserUuid,
		HashedPassword:    sql.NullString{String: "P4ssw0rd!", Valid: true},
		Email:             sql.NullString{String: util.RandomEmail(), Valid: true},
		FullName:          sql.NullString{String: util.RandomName(), Valid: true},
		PasswordChangedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.Email.String, updatedUser.Email)
	require.Equal(t, arg.UserUuid, updatedUser.UserUuid)
}

func TestUpdatePassword(t *testing.T) {
	user := GenerateUser(t)

	arg := UpdateUserPasswordParams{
		UserUuid:          user.UserUuid,
		HashedPassword:    sql.NullString{String: "P4ssw0rd!", Valid: true},
		PasswordChangedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	updatedUser, err := testQueries.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.UserUuid, updatedUser.UserUuid)
}

func GenerateUser(t *testing.T) CreateUserRow {
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
	require.NotEmpty(t, user)
	require.Equal(t, inputUser.Username, user.Username, "input and return username should be same")
	require.Equal(t, inputUser.Email, user.Email, "input and return email should be same")
	require.NotEmpty(t, user.UserUuid.String())
	require.NotEmpty(t, user.Role)

	return user
}
