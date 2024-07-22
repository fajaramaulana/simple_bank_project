package middleware

import (
	"fmt"
	"strings"

	"github.com/fajaramaulana/simple_bank_project/internal/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/token"
	"github.com/gin-gonic/gin"
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
