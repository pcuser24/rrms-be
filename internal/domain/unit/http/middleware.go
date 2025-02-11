package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth/http"
	unit_service "github.com/user2410/rrms-backend/internal/domain/unit/service"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

func CheckUnitManageability(s unit_service.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		tkPayload, ok := ctx.Locals(http.AuthorizationPayloadKey).(*token.Payload)
		if !ok {
			return ctx.SendStatus(fiber.StatusForbidden)
		}

		isManageable, err := s.CheckUnitManageability(puid, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isManageable {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this unit"})
		}

		return ctx.Next()
	}
}

const UnitIDLocalKey = "unitId"

func CheckUnitVisiblitiy(service unit_service.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		unitId, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		ctx.Locals(UnitIDLocalKey, unitId)

		var userId uuid.UUID = uuid.Nil
		tkPayload, ok := ctx.Locals(http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userId = tkPayload.UserID
		}

		isVisible, err := service.CheckVisibility(unitId, userId)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isVisible {
			return fiber.NewError(fiber.StatusForbidden, "this unit is not visible to you")
		}

		return ctx.Next()
	}
}
