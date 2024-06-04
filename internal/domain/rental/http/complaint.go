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

func (a *adapter) getRentalComplaintsOfUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		var query dto.GetRentalComplaintsOfUserQuery
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		res, err := a.service.GetRentalComplaintsOfUser(tkPayload.UserID, query)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "rental complaint not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) preCreateRentalComplaint() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.PreCreateRentalComplaint
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.service.PreCreateRentalComplaint(&payload, tkPayload.UserID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(payload)
	}
}

func (a *adapter) createRentalComplaint() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.CreateRentalComplaint
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.CreatorID = ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload).UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.CreateRentalComplaint(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getRentalComplaint() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rcId, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid rental complaint id: " + err.Error()})
		}

		rc, err := a.service.GetRentalComplaint(rcId)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "rental complaint not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(rc)
	}
}

func (a *adapter) getRentalComplaintsByRentalId() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rId, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid rental id: " + err.Error()})
		}

		res, err := a.service.GetRentalComplaintsByRentalId(rId)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) preCreateRentalComplaintReply() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.PreCreateRentalComplaint
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.service.PreCreateRentalComplaintReply(&payload, tkPayload.UserID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(payload)
	}
}

func (a *adapter) createRentalComplaintReply() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.CreateRentalComplaintReply
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.ReplierID = ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload).UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.CreateRentalComplaintReply(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getRentalComplaintReplies() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rcId, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid rental complaint id: " + err.Error()})
		}
		var query struct {
			Limit  int32 `query:"limit"`
			Offset int32 `query:"offset"`
		}
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		res, err := a.service.GetRentalComplaintReplies(rcId, query.Limit, query.Offset)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "rental complaint not found"})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) updateRentalComplaintStatus() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rcId, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid rental complaint id: " + err.Error()})
		}

		type Query struct {
			Status database.RENTALCOMPLAINTSTATUS `json:"status" validate:"required,oneof=RESOLVED CLOSED"`
		}
		var query Query
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		err = a.service.UpdateRentalComplaint(&dto.UpdateRentalComplaint{
			ID:     rcId,
			Status: query.Status,
			UserID: ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload).UserID,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
