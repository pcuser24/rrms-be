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
	propertyRoute := (*router).Group("/property")

	propertyRoute.Use(auth.NewAuthMiddleware(tokenMaker))

	propertyRoute.Post("/create", a.createProperty())
	propertyRoute.Get("/get-by-id/:id", checkPropertyOwnership(a.service), a.getPropertyById())
	propertyRoute.Patch("/update/:id", checkPropertyOwnership(a.service), a.updateProperty())
	propertyRoute.Delete("/delete/:id", checkPropertyOwnership(a.service), a.deleteProperty())

	// propertyRoute.Patch("/feature/:id", )
}

func checkPropertyOwnership(s Service) fiber.Handler {
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
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "operation not permitted on this property"})
		}

		ctx.Next()
		return nil
	}
}

func (a *adapter) createProperty() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateProperty
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		payload.OwnerID = tkPayload.UserID
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetFibValidationError(errs)})
			return nil
		}

		res, err := a.service.CreateProperty(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
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
		id := ctx.Params("id")
		puid, err := uuid.Parse(id)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		res, err := a.service.GetPropertyById(puid)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				responses.DBErrorResponse(ctx, pgErr)
				return nil
			}

			if errors.Is(err, sql.ErrNoRows) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("property with id=%s not found", puid.String())})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		if res.OwnerID != ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload).UserID {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "you are not authorized to view this property"})
		}
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
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetFibValidationError(errs)})
		}

		err := a.service.UpdateProperty(&payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
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
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return nil
	}
}
