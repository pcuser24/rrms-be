package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/domain/application"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const (
	ApplicationIdLocalKey = "aid"
)

func CheckApplicationVisibilty(s application.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tkPayload, ok := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		aid, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		c.Locals(ApplicationIdLocalKey, aid)

		isVisible, err := s.CheckApplicationVisibility(aid, tkPayload.UserID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isVisible {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You are not authorized to access this resource"})
		}

		return c.Next()
	}
}

func CheckApplicationUpdatability(s application.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tkPayload, ok := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		aid, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		c.Locals(ApplicationIdLocalKey, aid)

		isVisible, err := s.CheckApplicationUpdatability(aid, tkPayload.UserID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isVisible {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You are not authorized to access this resource"})
		}

		return c.Next()
	}
}
