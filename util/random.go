package util

import (
	"errors"

	"github.com/go-faker/faker/v4"
)

// Random Int
func RandomInt(n int) ([]int, error) {
	intRet, err := faker.RandomInt(n)
	if err != nil {
		return nil, errors.New("error generating random int")
	}

	return intRet, nil
}

func RandomName() string {
	nameRet := faker.Name()

	return nameRet
}

func RandomCurrency() string {
	return faker.Currency()
}
