package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/service"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

func (a *adapter) createRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.CreateRentalPayment
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.CreateRentalPayment(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rpId := ctx.Locals(RentalPaymentIDLocalKey).(int64)

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		rp, err := a.service.GetRentalPayment(rpId)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "rental payment not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		isVisible, err := a.service.CheckRentalVisibility(rp.RentalID, tkPayload.UserID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isVisible {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this rental payment"})
		}

		return ctx.Status(fiber.StatusOK).JSON(rp)
	}
}

func (a *adapter) getPaymentsOfRental() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rid := ctx.Locals(RentalIDLocalKey).(int64)

		res, err := a.service.GetPaymentsOfRental(rid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "rental not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getManagedRentalPayments() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		var query dto.GetManagedRentalPaymentsQuery
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.GetManagedRentalPayments(tkPayload.UserID, &query)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No payments found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) updatePlanRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalPaymentIDLocalKey).(int64)

		var payload dto.UpdatePlanRentalPayment
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if err := payload.Validate(); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.service.UpdateRentalPayment(id, tkPayload.UserID, &payload, database.RENTALPAYMENTSTATUSPLAN)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			if errors.Is(err, service.ErrInvalidPaymentTypeTransition) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) updateIssuedRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalPaymentIDLocalKey).(int64)

		var payload dto.UpdateIssuedRentalPayment
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.service.UpdateRentalPayment(id, tkPayload.UserID, &payload, database.RENTALPAYMENTSTATUSISSUED)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			if errors.Is(err, service.ErrInvalidPaymentTypeTransition) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) updatePendingRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalPaymentIDLocalKey).(int64)

		var payload dto.UpdatePendingRentalPayment
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.service.UpdateRentalPayment(
			id, tkPayload.UserID, &payload,
			utils.Ternary(payload.Status == database.RENTALPAYMENTSTATUSREQUEST2PAY,
				database.RENTALPAYMENTSTATUSPENDING,
				database.RENTALPAYMENTSTATUSREQUEST2PAY,
			),
		)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			if errors.Is(err, service.ErrInvalidPaymentTypeTransition) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) updatePartiallyPaidRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalPaymentIDLocalKey).(int64)

		var payload dto.UpdatePartiallyPaidRentalPayment
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.service.UpdateRentalPayment(
			id, tkPayload.UserID, &payload,
			database.RENTALPAYMENTSTATUSPARTIALLYPAID,
		)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			if errors.Is(err, service.ErrInvalidPaymentTypeTransition) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) updatePayfineRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Locals(RentalPaymentIDLocalKey).(int64)

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.service.UpdateRentalPayment(
			id, tkPayload.UserID, nil,
			database.RENTALPAYMENTSTATUSPAYFINE,
		)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			if errors.Is(err, service.ErrInvalidPaymentTypeTransition) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}
