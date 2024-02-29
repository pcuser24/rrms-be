package http

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/application"

	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service application.Service
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker) {
	applicationRoute := (*route).Group("/applications")

	applicationRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))

	applicationRoute.Post("/", a.createApplications())
	applicationRoute.Get("/my-applications",
		// TODO: A middleware to check if the user is a tenant
		a.getMyApplications(),
	)
	applicationRoute.Get("/to-me",
		// TODO: A middleware to check if the user is a property manager
		a.getApplicationsToMe(),
	)
	applicationRoute.Get("/application/:id",
		// TODO: A middleware to check if the user is a property manager
		a.getApplicationById(),
	)
	applicationRoute.Get("/ids", a.getApplicationsByIds())
	applicationRoute.Patch("/application/status/:id",
		// TODO: A middleware to check if the user is a property manager
		a.updateApplicationStatus(),
	)
	applicationRoute.Delete("/application/:id", a.deleteApplication())

}

func NewAdapter(service application.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) createApplications() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateApplication
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}
		payload.CreatorID = tkPayload.UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		// log.Println("CreateApplication payload:", payload)
		newApplication, err := a.service.CreateApplication(&payload)
		if err != nil {
			if dbErr, ok := err.(*database.TXError); ok {
				return responses.DBTXErrorResponse(ctx, dbErr)
			}
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			if errors.Is(err, application.ErrAlreadyApplied) ||
				errors.Is(err, application.ErrListingIsClosed) ||
				errors.Is(err, application.ErrInvalidApplicant) {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(newApplication)
	}
}

func (a *adapter) getMyApplications() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		query := new(dto.GetApplicationsToMeQuery)
		if err := query.QueryParser(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		validator := validation.GetDefaultValidator()
		validator.RegisterValidation(dto.PropertyFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, validation.GetValidationError(errs))
		}

		applications, err := a.service.GetApplicationsByUserId(tkPayload.UserID, query)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(applications)
	}
}

// Get applications to properties that the current user manages
func (a *adapter) getApplicationsToMe() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		query := new(dto.GetApplicationsToMeQuery)
		if err := query.QueryParser(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		validator := validation.GetDefaultValidator()
		validator.RegisterValidation(dto.PropertyFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		applications, err := a.service.GetApplicationsToUser(tkPayload.UserID, query)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(utils.Ternary(applications == nil, []model.ApplicationModel{}, applications))
	}
}

func (a *adapter) getApplicationById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		aid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		application, err := a.service.GetApplicationById(aid)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(application)
	}
}

func (a *adapter) getApplicationsByIds() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		query := new(dto.GetApplicationsByIdsQuery)
		if err := query.QueryParser(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}

		validator := validation.GetDefaultValidator()
		validator.RegisterValidation(dto.PropertyFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, validation.GetValidationError(errs))
		}

		var userId uuid.UUID = uuid.Nil
		tkPayload, ok := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userId = tkPayload.UserID
		}

		res, err := a.service.GetApplicationByIds(query.IDs, query.Fields, userId)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) updateApplicationStatus() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		aid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid application id"})
		}

		payload := new(dto.UpdateApplicationStatus)
		if err := ctx.BodyParser(payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		err = payload.Validate()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		err = a.service.UpdateApplicationStatus(aid, payload)
		if err != nil {
			if err == database.ErrRecordNotFound {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "application not found"})
			}
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (a *adapter) deleteApplication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		aid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		err = a.service.DeleteApplication((aid))
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return nil
	}
}
