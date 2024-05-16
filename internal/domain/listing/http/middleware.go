package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/listing/service"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const (
	ListingIDLocalKey = "listing_id"
)

func GetListingId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lid, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		c.Locals(ListingIDLocalKey, lid)

		return c.Next()
	}
}

func CheckListingManageability(s service.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		lid := ctx.Locals(ListingIDLocalKey).(uuid.UUID)
		tkPayload := ctx.Locals(http.AuthorizationPayloadKey).(*token.Payload)

		isManager, err := s.CheckListingOwnership(lid, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isManager {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this listing"})
		}

		ctx.Next()
		return nil
	}
}

func CheckListingVisibility(s service.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		lid := ctx.Locals(ListingIDLocalKey).(uuid.UUID)

		var userId uuid.UUID
		tkPayload, ok := ctx.Locals(http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userId = tkPayload.UserID
		}

		isVisible, err := s.CheckListingVisibility(lid, userId)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isVisible {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this listing"})
		}

		ctx.Next()
		return nil
	}
}
