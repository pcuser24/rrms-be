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
	unitRoute.Get("/search", a.searchUnits())
	unitRoute.Get("/ids", a.getUnitsByIds())

	unitRoute.Use(auth.AuthorizedMiddleware(tokenMaker))

	unitRoute.Post("/", a.createUnit())
	unitRoute.Patch("/unit/:id", CheckUnitManageability(a.uService), a.updateUnit())
	unitRoute.Delete("/unit/:id", CheckUnitManageability(a.uService), a.deleteUnit())
}

func (a *adapter) createUnit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.CreateUnit
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}
		if errs := utils.ValidateStruct(nil, payload); len(errs) > 0 && errs[0].Error {
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
		unitId := ctx.Locals(UnitIDLocalKey).(uuid.UUID)

		res, err := a.uService.GetUnitById(unitId)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return ctx.JSON(res)
	}
}

func (a *adapter) searchUnits() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var query dto.SearchUnitCombinationQuery
		if err := ctx.QueryParser(&query); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		if errs := utils.ValidateStruct(nil, query); len(errs) > 0 && errs[0].Error {
			return fiber.NewError(fiber.StatusBadRequest, utils.GetValidationError(errs))
		}

		res, err := a.uService.SearchUnit(&query)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.JSON(res)
	}
}

func (a *adapter) getUnitsByIds() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var query dto.GetUnitByIdsQuery
		if err := ctx.QueryParser(&query); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		validator := utils.GetDefaultValidator()
		validator.RegisterValidation(UnitFieldsLocalKey, dto.ValidateQuery)
		if errs := utils.ValidateStruct(validator, query); len(errs) > 0 && errs[0].Error {
			return fiber.NewError(fiber.StatusBadRequest, utils.GetValidationError(errs))
		}

		res, err := a.uService.GetUnitsByIds(query.IDs, query.Fields)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
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
		if errs := utils.ValidateStruct(nil, payload); len(errs) > 0 && errs[0].Error {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
			return nil
		}

		err := a.uService.UpdateUnit(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
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
