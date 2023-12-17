package property

import (
	"database/sql"
	"errors"
	"fmt"

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
	propertyRoute.Get("/property/:id", a.getPropertyById())

	propertyRoute.Use(auth.NewAuthMiddleware(tokenMaker))

	propertyRoute.Post("/", a.createProperty())
	propertyRoute.Patch("/property/:id", checkPropertyManageability(a.service), a.updateProperty())
	propertyRoute.Delete("/property/:id", checkPropertyManageability(a.service), a.deleteProperty())
}

func (a *adapter) createProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateProperty
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
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
		puid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		var userID uuid.UUID
		tkPayload, ok := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userID = tkPayload.UserID
		} else {
			userID = uuid.UUID{}
		}

		isVisible, err := a.service.CheckVisibility(puid, userID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isVisible {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "property is not visible to you"})
		}

		res, err := a.service.GetPropertyById(puid)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				return responses.DBErrorResponse(ctx, pgErr)
			}

			if errors.Is(err, sql.ErrNoRows) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("property with id=%s not found", puid.String())})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		// if res.CreatorID != ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload).UserID {
		// 	return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "you are not authorized to view this property"})
		// }
		return ctx.JSON(res)
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
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
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
