package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const (
	PropertyIDLocalKey = "property_id"
)

func GetPropertyId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		puid, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		c.Locals(PropertyIDLocalKey, puid)

		return c.Next()
	}
}

// Check whether the current user is a manager of the property
func CheckPropertyManageability(s property.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)
		tkPayload := ctx.Locals(http.AuthorizationPayloadKey).(*token.Payload)

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

// Check whether the current user is the owner of the property
func CheckPropertyOwnership(s property.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		ctx.Locals(PropertyIDLocalKey, puid)

		tkPayload, ok := ctx.Locals(http.AuthorizationPayloadKey).(*token.Payload)
		if !ok {
			return ctx.SendStatus(fiber.StatusForbidden)
		}

		isOwner, err := s.CheckOwnership(puid, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isOwner {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this property"})
		}

		return ctx.Next()
	}
}

// Check whether the property with given id  is visible or managed by the user of the current session
// should be stacked on top of AuthorizedMiddleware middleware
func CheckPropertyVisibility(service property.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		ctx.Locals(PropertyIDLocalKey, puid)

		var userID uuid.UUID = uuid.Nil
		tkPayload, ok := ctx.Locals(http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userID = tkPayload.UserID
		}

		isVisible, err := service.CheckVisibility(puid, userID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isVisible {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "this property is not visible to you"})
		}

		return ctx.Next()
	}
}
