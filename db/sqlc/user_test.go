package db

import (
	"context"
	"testing"
	"time"

	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	generateAccount(t)
}

func TestGetByEmail(t *testing.T) {
	user := GenerateUser(t)

	userByEmail, err := testStore.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, userByEmail)
	require.Equal(t, user.Email, userByEmail.Email)
	require.Equal(t, user.Username, userByEmail.Username)
	require.Equal(t, user.UserUuid, userByEmail.UserUuid)
}

func TestGetByUsername(t *testing.T) {
	user := GenerateUser(t)

	userByUsername, err := testStore.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userByUsername)
	require.Equal(t, user.Email, userByUsername.Email)
	require.Equal(t, user.Username, userByUsername.Username)
	require.Equal(t, user.UserUuid, userByUsername.UserUuid)
}

func TestUpdateUser(t *testing.T) {
	user := GenerateUser(t)
	hashedPassword, err := util.MakePasswordBcrypt("P4ssw0rdUpdate!")
	require.NoError(t, err)
	arg := UpdateUserParams{
		UserUuid:          user.UserUuid,
		HashedPassword:    pgtype.Text{String: hashedPassword, Valid: true},
		Email:             pgtype.Text{String: util.RandomEmail(), Valid: true},
		FullName:          pgtype.Text{String: util.RandomName(), Valid: true},
		PasswordChangedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	updatedUser, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.Email.String, updatedUser.Email)
	require.Equal(t, arg.UserUuid, updatedUser.UserUuid)
}

func TestUpdatePassword(t *testing.T) {
	user := GenerateUser(t)

	hashedPassword, err := util.MakePasswordBcrypt("P4ssw0rdUpdate!")

	require.NoError(t, err)

	arg := UpdateUserPasswordParams{
		UserUuid:          user.UserUuid,
		HashedPassword:    pgtype.Text{String: hashedPassword, Valid: true},
		PasswordChangedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	updatedUser, err := testStore.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.UserUuid, updatedUser.UserUuid)
}

func TestUpdateUserFullName(t *testing.T) {
	user := GenerateUser(t)

	res, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		UserUuid: user.UserUuid,
		FullName: pgtype.Text{String: "New Name", Valid: true},
	})

	require.NoError(t, err)
	require.NotEqual(t, user.FullName, res.FullName)
	require.Equal(t, "New Name", res.FullName)
	require.Equal(t, user.Email, res.Email)
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

	user, err := testStore.CreateUser(context.Background(), inputUser)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, inputUser.Username, user.Username, "input and return username should be same")
	require.Equal(t, inputUser.Email, user.Email, "input and return email should be same")
	require.NotEmpty(t, user.UserUuid.String())
	require.NotEmpty(t, user.Role)

	return user
}
