package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/rental"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const RentalIDLocalKey = "rental_id"

func CheckRentalVisibility(s rental.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message: Invalid rental id": err.Error()})
		}
		c.Locals(RentalIDLocalKey, rid)

		tkPayload := c.Locals(http.AuthorizationPayloadKey).(*token.Payload)

		isVisible, err := s.CheckRentalVisibility(rid, tkPayload.UserID)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if !isVisible {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this rental profile"})
		}

		return c.Next()
	}
}
