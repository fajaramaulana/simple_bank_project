package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const secretKeyMinimumLength = 32

type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new instance of JWTMaker with the provided secret key.
// It returns a Maker interface and an error if the secret key is invalid.
func NewJWTMaker(secretKey string) (Maker, error) {
	ErrInvalidSecretKey := fmt.Errorf("invalid secret key, must be at least %d characters", secretKeyMinimumLength)
	if len(secretKey) < secretKeyMinimumLength {
		return nil, ErrInvalidSecretKey
	}

	return &JWTMaker{secretKey}, nil
}

// CreateToken generates a new token for the given user UUID and duration.
// It returns the generated token as a string and any error encountered.
func (maker *JWTMaker) CreateToken(userUuid string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userUuid, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", payload, err
	}
	return token, payload, nil
}

// VerifyToken verifies the authenticity of the provided token.
// It returns the payload of the token if it is valid, or an error if the token is invalid.
// VerifyToken verifies the authenticity of a JWT token and returns the payload if the token is valid.
// It takes a token string as input and returns a pointer to the Payload struct and an error.
// If the token is invalid or an error occurs during verification, it returns nil and the corresponding error.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
