package token

import (
	"errors"
	"time"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/helper"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Role      string    `json:"role"`
}

func NewPayload(userUUIDString string, duration time.Duration, role string) (*Payload, error) {
	tokenUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// convert useruuid string to uuid.uuid
	userUUID, err := helper.ConvertStringToUUID(userUUIDString)
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenUUID,
		UserUUID:  userUUID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		Role:      role,
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}
