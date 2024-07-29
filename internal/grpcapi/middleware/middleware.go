package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/handler/token"
	"google.golang.org/grpc/metadata"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(ctx context.Context, tokenMaker token.Maker) (*token.Payload, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {

		authorizationHeader := md.Get(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			return nil, fmt.Errorf("authorization header is not provided")
		}

		authHeader := authorizationHeader[0]
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			return nil, fmt.Errorf("authorization header is not provided")
		}

		authType := strings.ToLower(fields[0])
		if authType != AuthorizationTypeBearer {
			return nil, fmt.Errorf("unsupported authorization type: %s", authType)
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return nil, err
		}

		return payload, nil
	}

	return nil, nil
}

func CheckRole(payload *token.Payload, role string) error {
	if payload.Role != role {
		return fmt.Errorf("access denied: role is not %s", role)
	}

	return nil
}
