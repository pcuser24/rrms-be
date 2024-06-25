package http

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	application_dto "github.com/user2410/rrms-backend/internal/domain/application/dto"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	auth_service "github.com/user2410/rrms-backend/internal/domain/auth/service"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_service "github.com/user2410/rrms-backend/internal/domain/property/service"
	rental_dto "github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker, authService auth_service.Service)
}

type adapter struct {
	service property_service.Service
}

func NewAdapter(service property_service.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker, authService auth_service.Service) {
	propertyRoute := (*router).Group("/properties")

	propertyRoute.Get("/property/:id",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		CheckPropertyVisibility(a.service),
		a.getPropertyById(),
	)
	propertyRoute.Get("/property/:id/units",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		CheckPropertyVisibility(a.service),
		a.getUnitsOfProperty(),
	)
	propertyRoute.Get("/ids",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		a.getPropertiesByIds(),
	)
	propertyRoute.Get("/new-manager-requests",
		auth_http.AuthorizedMiddleware(tokenMaker),
		a.getNewPropertyManagerRequests(),
	)
	propertyRoute.Get("/verification-status", a.getPropertiesVerificationStatus())

	propertyRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))

	propertyRoute.Post("/create/_pre", a.preCreateProperty())
	propertyRoute.Post("/create", a.createProperty())
	propertyRoute.Get("/managed-properties", a.getManagedProperties())
	propertyRoute.Group("/property/:id").Use(GetPropertyId())
	propertyRoute.Post("/property/:id/new-manager-requests",
		CheckPropertyManageability(a.service),
		a.createPropertyManagerRequest(),
	)
	propertyRoute.Patch("/property/:id/new-manager-requests/:requestId",
		a.updatePropertyManagerRequest(),
	)
	propertyRoute.Get("/property/:id/listings",
		CheckPropertyManageability(a.service),
		a.getListingsOfProperty(),
	)
	propertyRoute.Get("/property/:id/applications",
		CheckPropertyManageability(a.service),
		a.getApplicationsOfProperty(),
	)
	propertyRoute.Get("/property/:id/rentals",
		CheckPropertyManageability(a.service),
		a.getRentalsOfProperty(),
	)
	propertyRoute.Patch("/property/:id/_pre",
		CheckPropertyManageability(a.service),
		a.preUpdateProperty(),
	)
	propertyRoute.Patch("/property/:id",
		CheckPropertyManageability(a.service),
		a.updateProperty(),
	)
	propertyRoute.Delete("/property/:id",
		CheckPropertyManageability(a.service),
		a.deleteProperty(),
	)

	propertyVerificationRoute := propertyRoute.Group("/property/:id/verifications").Use(CheckPropertyManageability(a.service))
	propertyVerificationRoute.Post("/_pre", a.preCreatePropertyVerificationRequest())
	propertyVerificationRoute.Post("/", a.createVerificationRequest())
	propertyVerificationRoute.Get("/", a.getVerificationRequestsOfProperty())
	propertyVerificationRoute.Get("/verification/:vid", a.getVerificationRequest())

	propertyRoute.Get("/verifications", auth_http.AuthorizedMiddleware(tokenMaker), auth_http.AdminOnlyRoutes(authService), a.getVerificationRequests())
	propertyRoute.Patch("/verifications/:vid", auth_http.AuthorizedMiddleware(tokenMaker), auth_http.AdminOnlyRoutes(authService), a.updateVerificationRequestStatus())

}

func (a *adapter) preCreateProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.PreCreateProperty
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.service.PreCreateProperty(&payload, tkPayload.UserID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(payload)
	}
}

func (a *adapter) createProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateProperty
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.CreateProperty(&payload, tkPayload.UserID)
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

func (a *adapter) getPropertyById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)

		res, err := a.service.GetPropertyById(puid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getUnitsOfProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)

		res, err := a.service.GetUnitsOfProperty(puid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.JSON(res)
	}
}

func (a *adapter) getListingsOfProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)

		query := new(listing_dto.GetListingsOfPropertyQuery)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		validator := validator.New()
		validator.RegisterValidation(listing_dto.ListingFieldsLocalKey, listing_dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.GetListingsOfProperty(puid, query)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getApplicationsOfProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)

		query := new(application_dto.GetApplicationsOfPropertyQuery)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		validator := validation.GetDefaultValidator()
		validator.RegisterValidation(application_dto.ApplicationFieldsLocalKey, application_dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.GetApplicationsOfProperty(puid, query)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getRentalsOfProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)

		query := new(rental_dto.GetRentalsOfPropertyQuery)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		validator := validation.GetDefaultValidator()
		validator.RegisterValidation(rental_dto.RentalFieldsLocalKey, rental_dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.GetRentalsOfProperty(puid, query)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getPropertiesByIds() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		query := new(dto.GetPropertiesByIdsQuery)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		validator := validation.GetDefaultValidator()
		validator.RegisterValidation(dto.PropertyFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		var userId uuid.UUID = uuid.Nil
		tkPayload, ok := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userId = tkPayload.UserID
		}

		res, err := a.service.GetPropertiesByIds(query.IDs, query.Fields, userId)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getManagedProperties() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		query := new(dto.GetPropertiesQuery)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		validator := validator.New()
		validator.RegisterValidation("propertyFields", dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tokenPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		total, props, err := a.service.GetManagedProperties(tokenPayload.UserID, query)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"total": total,
			"items": props,
		})
	}
}

func (a *adapter) preUpdateProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.PreUpdateProperty
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.service.PreUpdateProperty(&payload, tkPayload.UserID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(payload)
	}
}

func (a *adapter) updateProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		puid, _ := uuid.Parse(id)

		var payload dto.UpdateProperty
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.ID = puid
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		err := a.service.UpdateProperty(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) deleteProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		puid, _ := uuid.Parse(id)

		err := a.service.DeleteProperty(puid)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}

func (a *adapter) createPropertyManagerRequest() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreatePropertyManagerRequest
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.CreatorID = tkPayload.UserID
		payload.PropertyID = puid
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.CreatePropertyManagerRequest(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			if dbErr, ok := err.(*database.TXError); ok {
				return responses.DBTXErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(res)
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getNewPropertyManagerRequests() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			limit = 100
		}
		offset, err := strconv.ParseInt(ctx.Query("offset"), 10, 64)
		if err != nil {
			offset = 0
		}

		res, err := a.service.GetNewPropertyManagerRequestsToUser(tkPayload.UserID, limit, offset)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no new requests"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) updatePropertyManagerRequest() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		pid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)
		requestId, err := strconv.ParseInt(ctx.Params("requestId"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request id"})
		}
		approved := (ctx.Query("approved") == "true")

		err = a.service.UpdatePropertyManagerRequest(pid, tkPayload.UserID, requestId, approved)
		if err != nil {
			if txErr, ok := err.(*database.TXError); ok {
				return responses.DBTXErrorResponse(ctx, txErr)
			}
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			if errors.Is(err, property_service.ErrUpdateRequestInfoMismatch) {
				return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
			}
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
