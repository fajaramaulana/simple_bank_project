package token

import "time"

// Maker is an interface that defines methods for creating and verifying tokens.
type Maker interface {
	// CreateToken generates a new token for the given user UUID and duration.
	// It returns the generated token as a string and any error encountered.
	CreateToken(userUuid string, duration time.Duration, role string) (string, *Payload, error)

	// VerifyToken verifies the authenticity of the provided token.
	// It returns the payload of the token if it is valid, or an error if the token is invalid.
	VerifyToken(token string) (*Payload, error)
}
