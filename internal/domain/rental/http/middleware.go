package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/rental/service"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const (
	RentalIDLocalKey          = "rental_id"
	RentalContractIDLocalKey  = "rental_contract_id"
	RentalPaymentIDLocalKey   = "rental_payment_id"
	RentalComplaintIDLocalKey = "rental_complaint_id"
)

func GetRentalID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message: Invalid rental id": err.Error()})
		}
		c.Locals(RentalIDLocalKey, id)

		return c.Next()
	}
}

func CheckRentalVisibility(s service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid := c.Locals(RentalIDLocalKey).(int64)

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

func GetContractID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message: Invalid rental id": err.Error()})
		}
		c.Locals(RentalContractIDLocalKey, id)

		return c.Next()
	}
}

func GetRentalPaymentID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message: Invalid rental id": err.Error()})
		}
		c.Locals(RentalPaymentIDLocalKey, id)

		return c.Next()
	}
}

func GetRentalComplaintID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message: Invalid rental id": err.Error()})
		}
		c.Locals(RentalComplaintIDLocalKey, id)

		return c.Next()
	}
}
