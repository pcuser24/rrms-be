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

	listingRoute.Post("/policy/add/:id", checkListingOwnership(a.lService), a.addListingPolicies())
	listingRoute.Delete("/policy/delete/:id", checkListingOwnership(a.lService), a.deleteListingPolicies())
	listingRoute.Post("/unit/add/:id", checkListingOwnership(a.lService), a.addListingUnits())
	listingRoute.Delete("/unit/delete/:id", checkListingOwnership(a.lService), a.deleteListingUnits())
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
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		// check ownership of target property
		isOwner, err := a.pService.CheckOwnership(payload.PropertyID, tkPayload.UserID)
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

		// validate units
		for _, pu := range payload.Units {
			isValid, err := a.uService.CheckUnitOfProperty(payload.PropertyID, pu.UnitID)
			if err != nil {
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}
			if !isValid {
				return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on unit id = " + pu.UnitID.String()})
			}
		}

		res, err := a.lService.CreateListing(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return ctx.Status(fiber.StatusCreated).JSON(res)
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

func (a *adapter) addListingPolicies() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		lid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []dto.CreateListingPolicy `json:"items" validate:"required,dive"`
		}
		if err := ctx.BodyParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		res, err := a.lService.AddListingPolicies(lid, query.Items)
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

func (a *adapter) deleteListingPolicies() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		lid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []int64 `json:"items" validate:"required,dive,gte=1"`
		}
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		err := a.lService.DeleteListingPolicies(lid, query.Items)
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

func (a *adapter) addListingUnits() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		lid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []dto.CreateListingUnit `json:"items" validate:"required,dive"`
		}
		if err := ctx.BodyParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		// validate input units
		for _, uid := range query.Items {
			isValid, err := a.lService.CheckValidUnitForListing(lid, uid.UnitID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			if !isValid {
				return fiber.NewError(fiber.StatusForbidden, "operation not permitted on unit id = "+uid.UnitID.String())
			}
		}

		res, err := a.lService.AddListingUnits(lid, query.Items)
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

func (a *adapter) deleteListingUnits() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		lid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []uuid.UUID `json:"items" validate:"required,dive,uuid4"`
		}
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		err := a.lService.DeleteListingUnits(lid, query.Items)
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
