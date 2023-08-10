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

	unitRoute.Patch("/amenity/update/:id", checkUnitOwnership(a.uService), a.updateUnitAmenities())
	unitRoute.Patch("/media/update/:id", checkUnitOwnership(a.uService), a.updateUnitMedium())
}

func checkUnitOwnership(s Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		puid, err := uuid.Parse(id)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		isOwner, err := s.CheckOwnership(puid, tkPayload.UserID)
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
		return ctx.JSON(res)
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

func (a *adapter) updateUnitMedium() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var payload struct {
			Items  []interface{} `json:"items" validate:"required"`
			Action string        `json:"action" validate:"required,oneof=add delete"`
		}
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
			return nil
		}

		var res interface{}
		var err error

		switch payload.Action {
		case "add":
			var items []dto.CreateUnitMedia
			for _, i := range payload.Items {
				var item dto.CreateUnitMedia
				err = utils.Map2JSONStruct(i, &item)
				if err != nil {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
				}
				items = append(items, item)
			}
			res, err = a.uService.AddUnitMedium(puid, items)
		case "delete":
			var items []int64
			for _, item := range payload.Items {
				i, ok := item.(float64)
				if !ok {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
				}
				items = append(items, int64(i))
			}
			err = a.uService.DeleteUnitMedium(puid, items)
		}
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		switch payload.Action {
		case "add":
			return ctx.Status(201).JSON(fiber.Map{"items": res})
		case "delete":
			return ctx.SendStatus(fiber.StatusNoContent)
		}

		return nil
	}
}

func (a *adapter) updateUnitAmenities() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var payload struct {
			Items  []interface{} `json:"items" validate:"required"`
			Action string        `json:"action" validate:"required,oneof=add delete replace"`
		}
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
			return nil
		}

		var res interface{}
		var err error

		switch payload.Action {
		case "add":
			var items []dto.CreateUnitAmenity
			for _, i := range payload.Items {
				var item dto.CreateUnitAmenity
				err = utils.Map2JSONStruct(i, &item)
				if err != nil {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
				}
				items = append(items, item)
			}
			res, err = a.uService.AddUnitAmenities(puid, items)
		case "delete":
			var items []int64
			for _, item := range payload.Items {
				i, ok := item.(float64)
				if !ok {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
				}
				items = append(items, int64(i))
			}
			err = a.uService.DeleteUnitAmenities(puid, items)
		}
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		switch payload.Action {
		case "add":
			return ctx.Status(201).JSON(fiber.Map{
				"items": res,
			})
		case "delete":
			return ctx.SendStatus(fiber.StatusNoContent)
		}

		return nil
	}
}
