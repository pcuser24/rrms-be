package listing

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
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
	listingRoute := (*router).Group("/listings")

	// non-required authorization routes
	listingRoute.Get("/", a.searchListings())
	listingRoute.Get("/listing/:id", a.getListingById())

	listingRoute.Use(auth.NewAuthMiddleware(tokenMaker))

	listingRoute.Post("/", a.createListing())
	listingRoute.Patch("/listing/:id", a.checkListingManageability(), a.updateListing())
	listingRoute.Delete("/listing/:id", a.checkListingManageability(), a.deleteListing())
}

func (a *adapter) createListing() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateListing
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.CreatorID = tkPayload.UserID
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}
		// log.Println("Passed struct validation")

		// check ownership of target property
		isManager, err := a.pService.CheckManageability(payload.PropertyID, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isManager {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this property"})
		}
		// log.Println("Passed property ownership check")

		// validate units
		for _, pu := range payload.Units {
			isValid, err := a.uService.CheckUnitOfProperty(payload.PropertyID, pu.UnitID)
			if err != nil {
				if dbErr, ok := err.(*pgconn.PgError); ok {
					return responses.DBErrorResponse(ctx, dbErr)
				}
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}
			if !isValid {
				return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on unit id = " + pu.UnitID.String()})
			}
		}
		// log.Println("Passed unit validation")

		res, err := a.lService.CreateListing(&payload)
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

func (a *adapter) searchListings() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.SearchListingCombinationQuery
		if err := ctx.QueryParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}
		log.Println(payload)

		res, err := a.lService.SearchListingCombination(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		return ctx.JSON(fiber.Map{
			"limit":  *payload.Limit,
			"offset": *payload.Offset,
			"count":  res.Count,
			"items":  res.Items,
		})
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
				return responses.DBErrorResponse(ctx, dbErr)
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
				return responses.DBErrorResponse(ctx, dbErr)
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
				return responses.DBErrorResponse(ctx, dbErr)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}
