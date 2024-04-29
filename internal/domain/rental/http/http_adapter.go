package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/rental"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service rental.Service
}

func NewAdapter(service rental.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker) {
	rentalRoute := (*route).Group("/rentals")
	rentalRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	rentalRoute.Post("/", a.createRental())
	rentalRoute.Get("/rental/:id", CheckRentalVisibility(a.service), a.getRental())
	rentalRoute.Patch("/rental/:id", CheckRentalVisibility(a.service), a.updateRental())
	rentalRoute.Get("/rental/:id/contract", CheckRentalVisibility(a.service), a.getRentalContract())
	rentalRoute.Get("/rental/:id/ping-contract", CheckRentalVisibility(a.service), a.pingContract())
	rentalRoute.Post("/rental/:id/contract", CheckRentalVisibility(a.service), a.createRentalContract())

	contractRoute := (*route).Group("/contracts")
	contractRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	contractRoute.Get("/contract/:id", a.getContract())
	// contractRoute.Patch("/contract/:id", a.updateContract())
	contractRoute.Patch("/contract/:id", a.updateContract())
	contractRoute.Patch("/contract/:id/content", a.updateContractContent())

	rentalPaymentRoute := (*route).Group("/rental-payments")
	rentalPaymentRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	rentalPaymentRoute.Post("/", a.createRentalPayment())
	rentalPaymentRoute.Get("/rental-payment/:id", a.getRentalPayment())
	rentalPaymentRoute.Get("/rental/:id", CheckRentalVisibility(a.service), a.getPaymentsOfRental())
	rentalPaymentRoute.Patch("/rental-payment/:id/plan", a.updatePlanRentalPayment())
	rentalPaymentRoute.Patch("/rental-payment/:id/issued", a.updateIssuedRentalPayment())
	rentalPaymentRoute.Patch("/rental-payment/:id/pending", a.updatePendingRentalPayment())
}

func (a *adapter) createRental() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateRental
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.CreatorID = tkPayload.UserID
		if err := payload.Validate(); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		res, err := a.service.CreateRental(&payload, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getRental() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rid := ctx.Locals(RentalIDLocalKey).(int64)

		res, err := a.service.GetRental(rid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "rental profile not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) updateRental() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		var payload dto.UpdateRental
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		err = a.service.UpdateRental(&payload, pid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (a *adapter) createRentalContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		var payload dto.CreateContract
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.RentalID = rid
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

func (a *adapter) getContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		res, err := a.service.GetContract(rid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) pingContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		res, err := a.service.PingRentalContract(rid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getRentalContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		res, err := a.service.GetRentalContract(rid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) updateContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		var payload dto.UpdateContract
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.ID = id
		payload.UserID = ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload).UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		err = a.service.UpdateContract(&payload)
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
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		var payload dto.UpdateContractContent
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.ID = id
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		err = a.service.UpdateContractContent(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}

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
		rpId, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid rental payment id: " + err.Error()})
		}

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

func (a *adapter) updatePlanRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		var payload dto.UpdatePlanRentalPayment
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err = a.service.UpdateRentalPayment(id, tkPayload.UserID, &payload, database.RENTALPAYMENTSTATUSPLAN)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			if errors.Is(err, rental.ErrInvalidTypeTransition) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) updateIssuedRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		var payload dto.UpdateIssuedRentalPayment
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err = a.service.UpdateRentalPayment(id, tkPayload.UserID, &payload, database.RENTALPAYMENTSTATUSISSUED)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			if errors.Is(err, rental.ErrInvalidTypeTransition) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) updatePendingRentalPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		var payload dto.UpdatePendingRentalPayment
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err = a.service.UpdateRentalPayment(
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

			if errors.Is(err, rental.ErrInvalidTypeTransition) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}
