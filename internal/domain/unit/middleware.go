package unit

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

func CheckUnitManageability(s Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		tkPayload, ok := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
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

func CheckUnitVisiblitiy(s Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		unitId, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		ctx.Locals(UnitIDLocalKey, unitId)

		var userId uuid.UUID = uuid.Nil
		tkPayload, ok := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userId = tkPayload.UserID
		}

		isVisible, err := s.CheckVisibility(unitId, userId)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if isVisible {
			return ctx.Next()
		} else {
			return fiber.NewError(fiber.StatusForbidden, "not allowed to get this unit")
		}
	}
}

const UnitFieldsLocalKey = "unitFields"

func GetUnitFieldsQuery() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var query dto.UnitFieldQuery
		if err := ctx.QueryParser(&query); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		validator := validator.New()
		validator.RegisterValidation(UnitFieldsLocalKey, dto.ValidateQuery)
		if errs := utils.ValidateStruct(validator, query); len(errs) > 0 && errs[0].Error {
			return fiber.NewError(fiber.StatusBadRequest, utils.GetValidationError(errs))
		}

		ctx.Locals(UnitFieldsLocalKey, &query)

		return ctx.Next()
	}
}
