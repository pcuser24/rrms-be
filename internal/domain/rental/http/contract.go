package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

func (a *adapter) createRentalContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		var payload dto.CreateContract
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.RentalID = id
		payload.UserID = ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload).UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.CreateContract(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getRentalContractsOfUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		query := new(dto.GetRentalContracts)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if err := query.ValidateQuery(); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		res, err := a.service.GetRentalContractsOfUser(tkPayload.UserID, query)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "contract not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalContractIDLocalKey).(int64)

		res, err := a.service.GetContract(id)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "contract not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) pingContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalIDLocalKey).(int64)

		res, err := a.service.PingRentalContract(id)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "contract not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getRentalContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalContractIDLocalKey).(int64)

		res, err := a.service.GetRentalContract(id)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "contract not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) updateContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalContractIDLocalKey).(int64)

		var payload dto.UpdateContract
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.ID = id
		payload.UserID = ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload).UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		err := a.service.UpdateContract(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (a *adapter) updateContractContent() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalContractIDLocalKey).(int64)

		var payload dto.UpdateContractContent
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.ID = id
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		err := a.service.UpdateContractContent(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
