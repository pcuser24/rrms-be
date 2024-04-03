package rental

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
	prentalRoute := (*route).Group("/rentals")
	prentalRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	prentalRoute.Post("/", a.createRental())
	prentalRoute.Get("/rental/:id", a.getRental())
	prentalRoute.Patch("/rental/:id", a.updateRental())
	prentalRoute.Get("/rental/:id/contract", a.getRentalContract())
	prentalRoute.Post("/rental/:id/contract", a.prepareRentalContract())
	prentalRoute.Patch("/rental/:id/contract", a.updateRentalContract())

	_ = (*route).Group("/rental")
}

func (a *adapter) createRental() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateRental
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.CreatorID = tkPayload.UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.CreateRental(&payload, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			if dbErr, ok := err.(*database.TXError); ok {
				return responses.DBTXErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getRental() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message: Invalid rental id": err.Error()})
		}

		res, err := a.service.GetRental(pid)
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
		// pid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		// if err != nil {
		// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		// }

		// res, err := a.service.GetRentalContract(pid)
		// if err != nil {
		// 	if errors.Is(err, database.ErrRecordNotFound) {
		// 		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
		// 	}

		// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		// }

		// return ctx.Status(fiber.StatusOK).JSON(res)
		return nil
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

func (a *adapter) updateRentalContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// pid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		// if err != nil {
		// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		// }

		// var payload dto.UpdateRentalContract
		// if err := ctx.BodyParser(&payload); err != nil {
		// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		// }
		// if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
		// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		// }

		// err = a.service.UpdateRentalContract(&payload, pid)
		// if err != nil {
		// 	if errors.Is(err, database.ErrRecordNotFound) {
		// 		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
		// 	}

		// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		// }

		// return ctx.SendStatus(fiber.StatusOK)
		return nil
	}
}

func (a *adapter) prepareRentalContract() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// pid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		// if err != nil {
		// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		// }

		// var payload dto.PrepareRentalContract
		// if err := ctx.BodyParser(&payload); err != nil {
		// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		// }
		// if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
		// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		// }

		// res, err := a.service.PrepareRentalContract(pid, &payload)
		// if err != nil {
		// 	if errors.Is(err, database.ErrRecordNotFound) {
		// 		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
		// 	}

		// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		// }

		// return ctx.Status(fiber.StatusOK).JSON(res)
		return nil
	}
}
