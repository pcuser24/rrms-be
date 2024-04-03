package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/listing"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

const (
	ListingIDLocalKey = "listing_id"
)

func CheckListingManageability(s listing.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		lid, err := uuid.Parse(id)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		ctx.Locals(ListingIDLocalKey, lid)

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

func CheckListingVisibility(s listing.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		lid, err := uuid.Parse(id)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		tkPayload := ctx.Locals(http.AuthorizationPayloadKey).(*token.Payload)

		isVisible, err := s.CheckListingVisibility(lid, tkPayload.UserID)
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
