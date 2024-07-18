package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func MakePasswordBcrypt(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password can't be empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordBcrypt(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
