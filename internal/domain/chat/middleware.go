package chat

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type GroupIdType = int64

const (
	GroupIDLocalKey = "groupId"
)

func CheckGroupMembership(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tkPayload, ok := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		groupId := c.Params("id")
		gid, err := strconv.ParseInt(groupId, 10, 64)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		s, err := s.CheckGroupMembership(tkPayload.UserID, gid)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if !s {
			return c.SendStatus(fiber.StatusForbidden)
		}
		c.Locals(GroupIDLocalKey, gid)
		return c.Next()
	}
}

const (
	AuthorizationHeaderKey  = "token"
	AuthorizationPayloadKey = "auth_payload"
)

// Middleware to check if the user is authorized to access the resource
func AuthorizedMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		accessToken := ctx.Query(AuthorizationHeaderKey)
		if len(accessToken) == 0 {
			return fiber.NewError(http.StatusUnauthorized, "token is not provided")
		}

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
