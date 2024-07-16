package util

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
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

func RandomEmail() string {
	return faker.Email()
}

func MakePasswordBcrypt(password string) string {
	// return faker.PasswordBcrypt(password)
	// generate password using bcrypt

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error generating password hash: %v", err)
	}

	return string(hashedPassword)
}
