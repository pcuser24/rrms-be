package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	auth_service "github.com/user2410/rrms-backend/internal/domain/auth/service"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationPayloadKey = "auth_payload"
	AuthorizationTypeBearer = "bearer"
)

// Middleware to check if the user is authorized to access the resource
func AuthorizedMiddleware(tokenMaker token.Maker) fiber.Handler {
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

// Middleware to get token payload if exists
// If the token is decoded successfully, the payload will be added to the context, whether it's valid or not
func GetAuthorizationMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.Get(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return ctx.Next()
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 || (strings.ToLower(fields[0]) != AuthorizationTypeBearer) {
			return ctx.Next()
		}

		accessToken := fields[1]
		payload, _ := tokenMaker.VerifyToken(accessToken)
		if payload == nil { // skip any logical error
			return ctx.Next()
		}

		ctx.Locals(AuthorizationPayloadKey, payload)

		return ctx.Next()
	}
}

func AddAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	userID uuid.UUID,
	duration time.Duration,
	options token.CreateTokenOptions,
) {
	token, payload, err := tokenMaker.CreateToken(userID, duration, options)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(AuthorizationHeaderKey, authorizationHeader)
}

func AdminOnlyRoutes(service auth_service.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload, ok := ctx.Locals(AuthorizationPayloadKey).(*token.Payload)
		if !ok {
			return ctx.SendStatus(fiber.StatusForbidden)
		}

		user, err := service.GetUserById(tkPayload.UserID)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		if user.Role != database.USERROLEADMIN {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted"})
		}

		return ctx.Next()
	}
}
