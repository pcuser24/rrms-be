package property

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const (
	PropertyIDLocalKey = "property_id"
)

func CheckPropertyManageability(s Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		tkPayload, ok := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		if !ok {
			return ctx.SendStatus(fiber.StatusForbidden)
		}

		isManager, err := s.CheckManageability(puid, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isManager {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this property"})
		}

		return ctx.Next()
	}
}

// should be stacked on top of AuthorizedMiddleware middleware
func CheckPropertyVisibility(service Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		ctx.Locals(PropertyIDLocalKey, puid)

		var userID uuid.UUID
		tkPayload, ok := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userID = tkPayload.UserID
		} else {
			userID = uuid.Nil
		}

		isVisible, err := service.CheckVisibility(puid, userID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isVisible {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "property is not visible to you"})
		}

		return ctx.Next()
	}
}