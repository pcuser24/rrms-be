package property

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service Service
}

func NewAdapter(service Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	propertyRoute := (*router).Group("/properties")

	propertyRoute.Get("/property/features", a.getAllFeatures())
	propertyRoute.Get("/property/:id", CheckPropertyVisibility(a.service), a.getPropertyById())
	propertyRoute.Get("/property/:id/units", CheckPropertyVisibility(a.service), a.getUnitsOfProperty())
	propertyRoute.Get("/ids", a.getPropertiesByIds())

	propertyRoute.Use(auth.AuthorizedMiddleware(tokenMaker))

	propertyRoute.Post("/", a.createProperty())
	propertyRoute.Get("/my-properties", a.getMyProperties())
	propertyRoute.Patch("/property/:id", CheckPropertyManageability(a.service), a.updateProperty())
	propertyRoute.Delete("/property/:id", CheckPropertyManageability(a.service), a.deleteProperty())
}

func (a *adapter) createProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateProperty
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(nil, payload); len(errs) > 0 {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
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

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		ctx.Status(fiber.StatusCreated).JSON(res)
		return nil
	}
}

func (a *adapter) getPropertyById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)

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
		puid := ctx.Locals(PropertyIDLocalKey).(uuid.UUID)

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
		validator := utils.GetDefaultValidator()
		validator.RegisterValidation(dto.PropertyFieldsLocalKey, dto.ValidateQuery)
		if errs := utils.ValidateStruct(validator, *query); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, utils.GetValidationError(errs))
		}

		res, err := a.service.GetPropertiesByIds(query.IDs, query.Fields)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		return ctx.JSON(res)
	}
}

func (a *adapter) getMyProperties() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		query := new(dto.GetPropertiesQuery)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		validator := validator.New()
		validator.RegisterValidation("propertyFields", dto.ValidateQuery)
		if errs := utils.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		tokenPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		res, err := a.service.GetPropertiesOfUser(tokenPayload.UserID, query.Fields)
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
		if errs := utils.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
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

func (a *adapter) getAllFeatures() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		res, err := a.service.GetAllFeatures()
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		return ctx.JSON(fiber.Map{
			"items": res,
		})
	}
}
