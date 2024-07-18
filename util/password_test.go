package util_test

import (
	"testing"

	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := util.RandomWord()

	hashedPassword1, err := util.MakePasswordBcrypt(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	// make test generate bycrpt error
	_, err = util.MakePasswordBcrypt("")
	require.Error(t, err)

	err = util.CheckPasswordBcrypt(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := util.RandomWord()
	err = util.CheckPasswordBcrypt(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := util.MakePasswordBcrypt(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
