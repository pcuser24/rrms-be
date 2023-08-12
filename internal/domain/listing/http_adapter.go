package listing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	pService property.Service
	uService unit.Service
	lService Service
}

func NewAdapter(lService Service, pService property.Service, uService unit.Service) Adapter {
	return &adapter{
		lService: lService,
		pService: pService,
		uService: uService,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	listingRoute := (*router).Group("/listing")

	listingRoute.Use(auth.NewAuthMiddleware(tokenMaker))

	listingRoute.Post("/create", a.createListing())
	listingRoute.Get("/get-by-id/:id", a.getListingById())
	listingRoute.Patch("/update/:id", checkListingOwnership(a.lService), a.updateListing())
	listingRoute.Delete("/delete/:id", checkListingOwnership(a.lService), a.deleteListing())

	listingRoute.Patch("/policy/update/:id", checkListingOwnership(a.lService), a.updateListingPolicy())
	listingRoute.Patch("/unit/update/:id", checkListingOwnership(a.lService), a.updateListingUnit())
}

func checkListingOwnership(s Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		lid, err := uuid.Parse(id)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		isOwner, err := s.CheckListingOwnership(lid, tkPayload.UserID)
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

func (a *adapter) createListing() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateListing
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}
		payload.CreatorID = tkPayload.UserID
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
			return nil
		}

		// check ownership of target property
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

		// validate units
		for _, unit := range payload.Units {
			isValid, err := a.lService.CheckValidUnitForListing(payload.PropertyID, unit.UnitID)
			if err != nil {
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}
			if !isValid {
				return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on unit id = " + unit.UnitID.String()})
			}
		}

		res, err := a.lService.CreateListing(&payload)
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

func (a *adapter) getListingById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		lid, err := uuid.Parse(id)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		res, err := a.lService.GetListingByID(lid)
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

func (a *adapter) updateListing() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		lid, _ := uuid.Parse(ctx.Params("id"))

		var payload dto.UpdateListing
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}
		payload.ID = lid
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
			return nil
		}

		err := a.lService.UpdateListing(&payload)
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

func (a *adapter) deleteListing() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		lid, err := uuid.Parse(id)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		err = a.lService.DeleteListing(lid)
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

func (a *adapter) updateListingPolicy() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		lid, _ := uuid.Parse(ctx.Params("id"))

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
			var items []dto.CreateListingPolicy
			for _, i := range payload.Items {
				var item dto.CreateListingPolicy
				err = utils.Map2JSONStruct(i, &item)
				if err != nil {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
				}
				items = append(items, item)
			}
			res, err = a.lService.AddListingPolicies(lid, items)
		case "delete":
			var items []int64
			for _, item := range payload.Items {
				i, ok := item.(float64)
				if !ok {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
				}
				items = append(items, int64(i))
			}
			err = a.lService.DeleteListingPolicies(lid, items)
		}
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
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

func (a *adapter) updateListingUnit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		lid, _ := uuid.Parse(ctx.Params("id"))

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
			var items []dto.CreateListingUnit
			for _, i := range payload.Items {
				var item dto.CreateListingUnit
				err = utils.Map2JSONStruct(i, &item)
				if err != nil {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
				}
				// validate units
				isValid, err := a.lService.CheckValidUnitForListing(lid, item.UnitID)
				if err != nil {
					return ctx.SendStatus(fiber.StatusInternalServerError)
				}
				if !isValid {
					return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on unit id = " + item.UnitID.String()})
				}
				items = append(items, item)
			}
			res, err = a.lService.AddListingUnits(lid, items)
		case "delete":
			var items []uuid.UUID
			for _, item := range payload.Items {
				i_str, ok := item.(string)
				if !ok {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
				}
				i, err := uuid.Parse(i_str)
				if err != nil {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
				}
				items = append(items, i)
			}
			err = a.lService.DeleteListingUnits(lid, items)
		}
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
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
