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

	(*router).Get("/property/amenities", a.getAllAmenities())
	(*router).Get("/property/features", a.getAllFeatures())

	propertyRoute = propertyRoute.Use(auth.NewAuthMiddleware(tokenMaker))

	propertyRoute.Post("/create", a.createProperty())
	propertyRoute.Get("/get-by-id/:id", checkPropertyOwnership(a.service), a.getPropertyById())
	propertyRoute.Patch("/update/:id", checkPropertyOwnership(a.service), a.updateProperty())
	propertyRoute.Delete("/delete/:id", checkPropertyOwnership(a.service), a.deleteProperty())

	propertyRoute.Post("/media/add/:id", checkPropertyOwnership(a.service), a.addPropertyMedium())
	propertyRoute.Delete("/media/delete/:id", checkPropertyOwnership(a.service), a.deletePropertyMedium())
	propertyRoute.Post("/amenity/add/:id", checkPropertyOwnership(a.service), a.addPropertyAmenities())
	propertyRoute.Delete("/amenity/delete/:id", checkPropertyOwnership(a.service), a.deletePropertyAmenities())
	propertyRoute.Post("/feature/add/:id", checkPropertyOwnership(a.service), a.addPropertyFeatures())
	propertyRoute.Delete("/feature/delete/:id", checkPropertyOwnership(a.service), a.deletePropertyFeatures())
	propertyRoute.Post("/tag/add/:id", checkPropertyOwnership(a.service), a.addPropertyTags())
	propertyRoute.Delete("/tag/delete/:id", checkPropertyOwnership(a.service), a.deletePropertyTags())
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
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
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
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
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

func (a *adapter) addPropertyMedium() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []dto.CreatePropertyMedia `json:"items" validate:"required,dive"`
		}
		if err := ctx.BodyParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		res, err := a.service.AddPropertyMedium(puid, query.Items)
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

func infoDeleteHandler(fn func(uuid.UUID, []int64) error) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []int64 `json:"items" validate:"required,dive,gte=1"`
		}
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		err := fn(puid, query.Items)
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

func (a *adapter) deletePropertyMedium() fiber.Handler {
	return infoDeleteHandler(a.service.DeletePropertyMedium)
}

func (a *adapter) addPropertyFeatures() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []dto.CreatePropertyFeature `json:"items" validate:"required,dive"`
		}
		if err := ctx.BodyParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		res, err := a.service.AddPropertyFeatures(puid, query.Items)
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

func (a *adapter) deletePropertyFeatures() fiber.Handler {
	return infoDeleteHandler(a.service.DeletePropertyFeatures)
}

func (a *adapter) addPropertyAmenities() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []dto.CreatePropertyAmenity `json:"items" validate:"required,dive"`
		}
		if err := ctx.BodyParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		res, err := a.service.AddPropertyAmenities(puid, query.Items)
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

func (a *adapter) deletePropertyAmenities() fiber.Handler {
	return infoDeleteHandler(a.service.DeletePropertyAmenities)
}

func (a *adapter) addPropertyTags() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []dto.CreatePropertyTag `json:"items" validate:"required,dive"`
		}
		if err := ctx.BodyParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		res, err := a.service.AddPropertyTags(puid, query.Items)
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

func (a *adapter) deletePropertyTags() fiber.Handler {
	return infoDeleteHandler(a.service.DeletePropertyTags)
}

func (a *adapter) getAllAmenities() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		res, err := a.service.GetAllAmenities()
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		return ctx.JSON(fiber.Map{
			"items": res,
		})
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
