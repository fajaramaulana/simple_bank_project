package util

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

func RandomUsername() string {
	return faker.Username()
}
func RandomCurrency() string {
	// random only USD EUR IDR
	currency := []string{"USD", "EUR", "IDR"}
	return currency[rand.Intn(len(currency))]
}

func RandomMoney(r *rand.Rand, min, max float64) string {
	amount := min + r.Float64()*(max-min)
	return fmt.Sprintf("%.2f", amount)
}

func RandomEmail() string {
	return faker.Email()
}

func RandomWord() string {
	return faker.Word()
}

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandomRole() string {
	// random only USD EUR IDR
	role := []string{"cusomer", "admin", "superadmin"}
	return role[rand.Intn(len(role))]
}

func RandomUUID() uuid.UUID {
	return uuid.New()
}
