package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service property.Service
}

func NewAdapter(service property.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	propertyRoute := (*router).Group("/properties")

	propertyRoute.Get("/property/:id",
		auth.GetAuthorizationMiddleware(tokenMaker),
		property.CheckPropertyVisibility(a.service),
		a.getPropertyById(),
	)
	propertyRoute.Get("/property/:id/units",
		property.CheckPropertyVisibility(a.service),
		a.getUnitsOfProperty(),
	)
	propertyRoute.Get("/ids",
		auth.GetAuthorizationMiddleware(tokenMaker),
		a.getPropertiesByIds(),
	)

	propertyRoute.Use(auth.AuthorizedMiddleware(tokenMaker))

	propertyRoute.Post("/", a.createProperty())
	propertyRoute.Get("/my-properties", a.getManagedProperties())
	propertyRoute.Patch("/property/:id",
		property.CheckPropertyManageability(a.service),
		a.updateProperty(),
	)
	propertyRoute.Delete("/property/:id",
		property.CheckPropertyManageability(a.service),
		a.deleteProperty(),
	)
}

func (a *adapter) createProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateProperty
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
			return nil
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
		puid := ctx.Locals(property.PropertyIDLocalKey).(uuid.UUID)

		res, err := a.service.GetPropertyById(puid)
		if err != nil {
			if err == database.ErrRecordNotFound {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.JSON(res)
	}
}

func (a *adapter) getUnitsOfProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid := ctx.Locals(property.PropertyIDLocalKey).(uuid.UUID)

		res, err := a.service.GetUnitsOfProperty(puid)
		if err != nil {
			if err == database.ErrRecordNotFound {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "property not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.JSON(res)
	}
}

func (a *adapter) getPropertiesByIds() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		query := new(dto.GetPropertiesByIdsQuery)
		if err := query.QueryParser(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}

		validator := validation.GetDefaultValidator()
		validator.RegisterValidation(dto.PropertyFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, validation.GetValidationError(errs))
		}

		var userId uuid.UUID
		tkPayload, ok := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userId = tkPayload.UserID
		} else {
			userId = uuid.Nil
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

		tokenPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		res, err := a.service.GetManagedProperties(tokenPayload.UserID, query.Fields)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) updateProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		puid, _ := uuid.Parse(id)

		var payload dto.UpdateProperty
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
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
