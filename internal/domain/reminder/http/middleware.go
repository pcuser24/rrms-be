package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/reminder"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const ReminderIdLocalKey = "reminder_id"

func GetReminderId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		c.Locals(ReminderIdLocalKey, id)
		return c.Next()
	}
}
func CheckReminderVisibility(s reminder.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Locals(ReminderIdLocalKey).(int64)
		tkPayload := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		visible, err := s.CheckReminderVisibility(id, tkPayload.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !visible {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You are not allowed to access this reminder"})
		}
		return c.Next()
	}
}
