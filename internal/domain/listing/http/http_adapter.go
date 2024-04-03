package http

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/listing"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/types"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	pService property.Service
	uService unit.Service
	lService listing.Service
}

func NewAdapter(lService listing.Service, pService property.Service, uService unit.Service) Adapter {
	return &adapter{
		lService: lService,
		pService: pService,
		uService: uService,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	listingRoute := (*router).Group("/listings")

	listingRoute.Get("/", a.searchListings())
	listingRoute.Get("/listing/:id",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		CheckListingVisibility(a.lService),
		a.getListingById(),
	)
	listingRoute.Get("/ids", a.getListingsByIds())
	listingRoute.Get("/listing/:id/application-link", a.verifyApplicationLink())

	listingRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))

	listingRoute.Post("/", a.createListing())
	listingRoute.Get("/my-listings", a.getMyListings())
	listingRoute.Post("/listing/:id/payment", CheckListingManageability(a.lService), a.createListingPayment())
	listingRoute.Post("/listing/:id/application-link", CheckListingManageability(a.lService), a.createApplicationLink())
	listingRoute.Patch("/listing/:id", CheckListingManageability(a.lService), a.updateListing())
	listingRoute.Delete("/listing/:id", CheckListingManageability(a.lService), a.deleteListing())
}

func (a *adapter) createListing() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateListing
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.CreatorID = tkPayload.UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}
		// log.Println("Passed struct validation")

		// check managability of the target property
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

		res, err := a.lService.CreateListing(&payload)
		if err != nil {
			if dbErr, ok := err.(*database.TXError); ok {
				return responses.DBTXErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		// res, err := a.lService.CreateListing(&payload)
		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) searchListings() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(dto.SearchListingCombinationQuery)
		if err := payload.QueryParser(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, *payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}
		// log.Println(payload)

		// make sure properties are public and listings not expired
		payload.PIsPublic = types.Ptr[bool](true)
		payload.LMinExpiredAt = types.Ptr[time.Time](time.Now())

		res, err := a.lService.SearchListingCombination(payload)
		if err != nil {
			if err == database.ErrRecordNotFound {
				return ctx.SendStatus(fiber.StatusNotFound)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getMyListings() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		query := new(dto.GetListingsQuery)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		validator := validator.New()
		validator.RegisterValidation("listingFields", dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tokenPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		res, err := a.lService.GetListingsOfUser(tokenPayload.UserID, query.Fields)
		if err != nil {
			if err == database.ErrRecordNotFound {
				return ctx.SendStatus(fiber.StatusNotFound)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
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
			if err == database.ErrRecordNotFound {
				return ctx.SendStatus(fiber.StatusNotFound)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}
		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getListingsByIds() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		query := new(dto.GetListingsByIdsQuery)
		if err := query.QueryParser(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		validator := validation.GetDefaultValidator()
		validator.RegisterValidation(dto.ListingFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, validation.GetValidationError(errs))
		}

		res, err := a.lService.GetListingsByIds(query.IDs, query.Fields)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
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
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
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

func (a *adapter) createListingPayment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.CreateListingPayment
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		payload.ListingId = ctx.Locals(ListingIDLocalKey).(uuid.UUID)
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		payload.UserId = tkPayload.UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.lService.CreateListingPayment(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) createApplicationLink() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.CreateApplicationLink
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.ListingId = ctx.Locals(ListingIDLocalKey).(uuid.UUID)
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		query, err := a.lService.CreateApplicationLink(&payload)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"url": fmt.Sprintf("%s/application/%s?%s", ctx.Get("Origin"), payload.ListingId.String(), query)})
	}
}

func (a *adapter) verifyApplicationLink() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var query dto.VerifyApplicationLink
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		id, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid id"})
		}
		query.ListingId = id
		if errs := validation.ValidateStruct(nil, query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.lService.VerifyApplicationLink(&query)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if !res {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "invalid application link"})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
