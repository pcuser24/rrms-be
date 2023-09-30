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

	(*router).Get("/property/features", a.getAllFeatures())
	(*router).Get("/property/get-by-id/:id", a.getPropertyById())

	propertyRoute = propertyRoute.Use(auth.NewAuthMiddleware(tokenMaker))

	propertyRoute.Post("/create", a.createProperty())
	propertyRoute.Patch("/update/:id", checkPropertyManageability(a.service), a.updateProperty())
	propertyRoute.Delete("/delete/:id", checkPropertyManageability(a.service), a.deleteProperty())

	propertyRoute.Post("/media/add/:id", checkPropertyManageability(a.service), a.addPropertyMedia())
	propertyRoute.Delete("/media/delete/:id", checkPropertyManageability(a.service), a.deletePropertyMedia())
	propertyRoute.Post("/feature/add/:id", checkPropertyManageability(a.service), a.addPropertyFeatures())
	propertyRoute.Delete("/feature/delete/:id", checkPropertyManageability(a.service), a.deletePropertyFeatures())
	propertyRoute.Post("/tag/add/:id", checkPropertyManageability(a.service), a.addPropertyTags())
	propertyRoute.Delete("/tag/delete/:id", checkPropertyManageability(a.service), a.deletePropertyTags())
	propertyRoute.Post("/manager/add/:id", checkPropertyManageability(a.service), a.addPropertyManagers())
	propertyRoute.Delete("/manager/delete/:id", checkPropertyManageability(a.service), a.deletePropertyManagers())
}

func checkPropertyManageability(s Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		isManager, err := s.CheckManageability(puid, tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		if !isManager {
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
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
			return nil
		}

		res, err := a.service.CreateProperty(&payload, tkPayload.UserID)
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
		puid, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		var userID uuid.UUID
		tkPayload, ok := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userID = tkPayload.UserID
		}

		isVisible, err := a.service.CheckVisibility(puid, userID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
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
				responses.DBErrorResponse(ctx, pgErr)
				return nil
			}

			if errors.Is(err, sql.ErrNoRows) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("property with id=%s not found", puid.String())})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		if res.CreatorID != ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload).UserID {
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

func (a *adapter) addPropertyMedia() fiber.Handler {
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

		res, err := a.service.AddPropertyMedia(puid, query.Items)
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

func (a *adapter) deletePropertyMedia() fiber.Handler {
	return infoDeleteHandler(a.service.DeletePropertyMedia)
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

func (a *adapter) addPropertyManagers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))

		var query struct {
			Items []dto.CreatePropertyManager `json:"items" validate:"required,dive"`
		}
		if err := ctx.BodyParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := utils.ValidateStruct(query); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetValidationError(errs)})
		}

		res, err := a.service.AddPropertyManagers(puid, query.Items)
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

func (a *adapter) deletePropertyManagers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		puid, _ := uuid.Parse(ctx.Params("id"))
		mid, err := uuid.Parse(ctx.Query("mid"))
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid manager id"})
		}

		err = a.service.DeletePropertyManager(puid, mid)
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
