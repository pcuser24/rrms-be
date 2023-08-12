package unit

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
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
	unitRoute := (*router).Group("/unit")

	(*router).Get("/unit/amenities", a.getAllAmenities())

	unitRoute.Use(auth.NewAuthMiddleware(tokenMaker))

	unitRoute.Post("/create", a.createUnit())
	unitRoute.Get("/get-by-id/:id", checkUnitOwnership(a.uService), a.getUnitById())
	unitRoute.Get("/get-by-property-id/:id", a.getUnitsOfProperty())
	unitRoute.Patch("/update/:id", checkUnitOwnership(a.uService), a.updateUnit())
	unitRoute.Delete("/delete/:id", checkUnitOwnership(a.uService), a.deleteUnit())

	unitRoute.Post("/amenity/add/:id", checkUnitOwnership(a.uService), a.addUnitAmenities())
	unitRoute.Delete("/amenity/delete/:id", checkUnitOwnership(a.uService), a.deleteUnitAmenities())
	unitRoute.Post("/media/add/:id", checkUnitOwnership(a.uService), a.addUnitMedium())
	unitRoute.Delete("/media/delete/:id", checkUnitOwnership(a.uService), a.deleteUnitMedium())
}

func checkUnitOwnership(s Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		puid, err := uuid.Parse(id)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		isOwner, err := s.CheckUnitOwnership(puid, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isOwner {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this unit"})
		}

		ctx.Next()
		return nil
	}
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
		isOwner, err := a.pService.CheckOwnership(payload.PropertyID, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		if !isOwner {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this property"})
		}

		res, err := a.uService.CreateUnit(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getUnitById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		res, err := a.uService.GetUnitById(uid)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
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
		isOwner, err := a.pService.CheckOwnership(pid, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isOwner {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this property"})
		}

		res, err := a.uService.GetUnitsOfProperty(pid)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
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
				responses.DBErrorResponse(ctx, dbErr)
				return nil
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
				responses.DBErrorResponse(ctx, dbErr)
				return nil
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

func infoDeleteHandler(fn func(uuid.UUID, []int64) error) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []int64 `json:"items" validate:"required,dive,gte=1"`
		}
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		err := fn(puid, query.Items)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusNoContent)
	}
}

func (a *adapter) addUnitMedium() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []dto.CreateUnitMedia `json:"items" validate:"required,dive"`
		}
		if err := ctx.BodyParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		res, err := a.uService.AddUnitMedium(puid, query.Items)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(201).JSON(fiber.Map{"items": res})
	}
}

func (a *adapter) deleteUnitMedium() fiber.Handler {
	return infoDeleteHandler(a.uService.DeleteUnitMedium)
}

func (a *adapter) addUnitAmenities() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []dto.CreateUnitAmenity `json:"items" validate:"required,dive"`
		}
		if err := ctx.BodyParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		res, err := a.uService.AddUnitAmenities(puid, query.Items)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(201).JSON(fiber.Map{"items": res})
	}
}

func (a *adapter) deleteUnitAmenities() fiber.Handler {
	return infoDeleteHandler(a.uService.DeleteUnitAmenities)
}
