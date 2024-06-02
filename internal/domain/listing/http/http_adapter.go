package http

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	"github.com/user2410/rrms-backend/internal/domain/listing/utils"
	property_service "github.com/user2410/rrms-backend/internal/domain/property/service"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	pService property_service.Service
	uService unit.Service
	lService listing_service.Service
}

func NewAdapter(lService listing_service.Service, pService property_service.Service, uService unit.Service) Adapter {
	return &adapter{
		lService: lService,
		pService: pService,
		uService: uService,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	listingRoute := (*router).Group("/listings")

	listingRoute.Get("/search",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		a.searchListings(),
	)
	listingRoute.Get("/ids", auth_http.GetAuthorizationMiddleware(tokenMaker), a.getListingsByIds())
	listingRoute.Get("/listing/:id/application-link", a.verifyApplicationLink())
	listingRoute.Get("/listing/:id",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		GetListingId(),
		CheckListingVisibility(a.lService),
		a.getListingById(),
	)

	listingRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))

	listingRoute.Post("/", a.createListing())
	listingRoute.Get("/managed-listings", a.getManagedListings())

	listingRoute.Group("/listing/:id").Use(GetListingId())
	listingRoute.Post("/listing/:id/application-link", CheckListingManageability(a.lService), a.createApplicationLink())
	listingRoute.Patch("/listing/:id", CheckListingManageability(a.lService), a.updateListing())
	listingRoute.Get("/listing/:id/payments", CheckListingManageability(a.lService), a.getListingPayments())
	listingRoute.Patch("/listing/:id/upgrade", CheckListingManageability(a.lService), a.upgradeListing())
	listingRoute.Patch("/listing/:id/extend", CheckListingManageability(a.lService), a.extendListing())
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

		var userId uuid.UUID
		tkPayload, ok := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userId = tkPayload.UserID
		}

		res, err := a.lService.SearchListingCombination(payload, userId)
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

func (a *adapter) getManagedListings() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		query := new(dto.GetListingsQuery)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		validator := validator.New()
		validator.RegisterValidation(dto.ListingFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tokenPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		count, listings, err := a.lService.GetListingsOfUser(tokenPayload.UserID, query)
		if err != nil {
			if err == database.ErrRecordNotFound {
				return ctx.SendStatus(fiber.StatusNotFound)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": count,
			"items": listings,
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

		uid := uuid.Nil
		tkPayload, ok := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			tkPayload.UserID = uid
		}
		res, err := a.lService.GetListingsByIds(uid, query.IDs, query.Fields)
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
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		err := a.lService.UpdateListing(lid, &payload)
		if err != nil {
			if txErr, ok := err.(*database.TXError); ok {
				if txErr.RollbackErr != nil || txErr.CommitErr != nil {
					return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": txErr.RollbackErr.Error()})
				}
				if dbErr, ok := txErr.Err.(*pgconn.PgError); ok {
					return responses.DBErrorResponse(ctx, dbErr)
				}
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) deleteListing() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := a.lService.DeleteListing(ctx.Locals(ListingIDLocalKey).(uuid.UUID))
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
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
		query.ListingId = ctx.Locals(ListingIDLocalKey).(uuid.UUID)
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

// Request to upgrade listing
func (a *adapter) upgradeListing() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload struct {
			Priority int `json:"priority" validate:"required,gt=0,lte=4"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		payment, err := a.lService.UpgradeListing(
			c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload).UserID,
			c.Locals(ListingIDLocalKey).(uuid.UUID),
			payload.Priority,
		)
		if err != nil {
			if errors.Is(err, utils.ErrInvalidPriority) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(c, dbErr)
			}

			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(payment)
	}
}

func (a *adapter) extendListing() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload struct {
			Duration int `json:"priority" validate:"required,gt=0"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		payment, err := a.lService.ExtendListing(
			c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload).UserID,
			c.Locals(ListingIDLocalKey).(uuid.UUID),
			payload.Duration,
		)
		if err != nil {
			if errors.Is(err, utils.ErrInvalidDuration) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			}
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(c, dbErr)
			}

			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(payment)
	}
}

func (a *adapter) getListingPayments() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		res, err := a.lService.GetListingPayments(ctx.Locals(ListingIDLocalKey).(uuid.UUID))
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
