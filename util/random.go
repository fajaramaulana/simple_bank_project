package util

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
)

func NewRandomMoneyGenerator() *rand.Rand {
	seed := rand.NewSource(time.Now().UnixNano())
	return rand.New(seed)
}

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

func RandomMoney(r *rand.Rand, min, max float64) string {
	amount := min + r.Float64()*(max-min)
	return fmt.Sprintf("%.2f", amount)
}
