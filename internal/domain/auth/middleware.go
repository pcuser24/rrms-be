package auth

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationPayloadKey = "auth_payload"
	AuthorizationTypeBearer = "bearer"
)

func NewAuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.Get(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return fiber.NewError(http.StatusUnauthorized, "authorization header is not provided")
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 || (strings.ToLower(fields[0]) != AuthorizationTypeBearer) {
			return fiber.NewError(http.StatusUnauthorized, "invalid authorization header")
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}

		if payload.TokenType != token.AccessToken {
			return fiber.NewError(http.StatusUnauthorized, "invalid token type")
		}

		ctx.Locals(AuthorizationPayloadKey, payload)
		ctx.Next()

		return nil
	}
}
