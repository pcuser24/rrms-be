package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	application "github.com/user2410/rrms-backend/internal/domain/application/service"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const (
	ApplicationIdLocalKey = "aid"
)

func GetApplicationId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		aid, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		c.Locals(ApplicationIdLocalKey, aid)

		return c.Next()
	}
}

func CheckApplicationVisibilty(s application.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tkPayload := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		aid := c.Locals(ApplicationIdLocalKey).(int64)

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
		tkPayload := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		aid := c.Locals(ApplicationIdLocalKey).(int64)

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
