package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := "authorization header is not provided"
			helper.ReturnJSONAbort(c, 401, err, nil)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := "authorization header is not provided"
			helper.ReturnJSONAbort(c, 401, err, nil)
			return
		}

		authorizationType := strings.ToLower(fields[0])

		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			helper.ReturnJSONAbort(c, 401, err.Error(), nil)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			helper.ReturnJSONAbort(c, 401, err.Error(), nil)
			return
		}

		c.Set(AuthorizationPayloadKey, payload)
		c.Next()
	}
}

func AddAuthorizationTest(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	duration time.Duration,
	role string,
) {
	uuidToken, err := uuid.NewRandom()
	require.NoError(t, err)

	uuidTokenString := uuidToken.String()

	token, payload, err := tokenMaker.CreateToken(uuidTokenString, duration, role)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(AuthorizationHeaderKey, authorizationHeader)
}

func AddAuthorizationTestAPI(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	uuidTokenString string,
	duration time.Duration,
	role string,
) {

	token, payload, err := tokenMaker.CreateToken(uuidTokenString, duration, role)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(AuthorizationHeaderKey, authorizationHeader)
}
