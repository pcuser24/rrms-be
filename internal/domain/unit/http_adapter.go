package unit

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	uService Service
	pService property.Service
}

func NewAdapter(uService Service, pService property.Service) Adapter {
	return &adapter{
		uService: uService,
		pService: pService,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	unitRoute := (*router).Group("/units")

	unitRoute.Get("/unit/amenities", a.getAllAmenities())
	unitRoute.Get("/unit/:id", a.getUnitById())

	unitRoute.Use(auth.NewAuthMiddleware(tokenMaker))

	unitRoute.Post("/", a.createUnit())
	unitRoute.Get("/property/:id", a.getUnitsOfProperty())
	unitRoute.Patch("/unit/:id", checkUnitManageability(a.uService), a.updateUnit())
	unitRoute.Delete("/unit/:id", checkUnitManageability(a.uService), a.deleteUnit())
}

func (a *adapter) createUnit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.CreateUnit
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
			return nil
		}

		// check ownership of target property
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		isManageable, err := a.pService.CheckManageability(payload.PropertyID, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		if !isManageable {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this property"})
		}

		res, err := a.uService.CreateUnit(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			if dbErr, ok := err.(*database.TXError); ok {
				return responses.DBTXErrorResponse(ctx, dbErr)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getUnitById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		uid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		var userID uuid.UUID
		tkPayload, ok := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userID = tkPayload.UserID
		}

		isVisible, err := a.uService.CheckVisibility(uid, userID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isVisible {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "property is not visible to you"})
		}

		res, err := a.uService.GetUnitById(uid)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		return ctx.JSON(res)
	}
}

func (a *adapter) getUnitsOfProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		pid, err := uuid.Parse(id)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		isManageable, err := a.pService.CheckManageability(pid, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isManageable {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this property"})
		}

		res, err := a.uService.GetUnitsOfProperty(pid)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return ctx.JSON(fiber.Map{
			"items": res,
		})
	}
}

func (a *adapter) updateUnit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		uid, _ := uuid.Parse(ctx.Params("id"))

		var payload dto.UpdateUnit
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}
		payload.ID = uid
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
			return nil
		}

		err := a.uService.UpdateUnit(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) deleteUnit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		err = a.uService.DeleteUnit(uid)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) getAllAmenities() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		res, err := a.uService.GetAllAmenities()
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		return ctx.JSON(fiber.Map{
			"items": res,
		})
	}
}
