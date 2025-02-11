package http

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	application_service "github.com/user2410/rrms-backend/internal/domain/application/service"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"

	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	lService listing_service.Service
	aService application_service.Service
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker) {
	applicationRoute := (*route).Group("/applications")

	// applicationRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))

	applicationRoute.Post("/create/_pre",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		a.preCreateApplications(),
	)
	applicationRoute.Post("/create",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		a.createApplications(),
	)
	applicationRoute.Get("/my-applications",
		auth_http.AuthorizedMiddleware(tokenMaker),
		a.getMyApplications(),
	)
	applicationRoute.Get("/to-me",
		auth_http.AuthorizedMiddleware(tokenMaker),
		a.getApplicationsToMe(),
	)
	applicationRoute.Get("/ids",
		auth_http.AuthorizedMiddleware(tokenMaker),
		a.getApplicationsByIds(),
	)
	applicationRoute.Group("/application/:id").Use(GetApplicationId())
	applicationRoute.Get("/application/:id",
		auth_http.AuthorizedMiddleware(tokenMaker),
		CheckApplicationVisibilty(a.aService),
		a.getApplicationById(),
	)
	applicationRoute.Patch("/application/:id/status",
		auth_http.AuthorizedMiddleware(tokenMaker),
		CheckApplicationUpdatability(a.aService),
		a.updateApplicationStatus(),
	)
	applicationRoute.Post("/application/:id/msg-group",
		auth_http.AuthorizedMiddleware(tokenMaker),
		CheckApplicationVisibilty(a.aService),
		a.createApplicationMsgGroup(),
	)
	applicationRoute.Get("/application/:id/msg-group",
		auth_http.AuthorizedMiddleware(tokenMaker),
		CheckApplicationVisibilty(a.aService),
		a.getApplicationMsgGroup(),
	)
	applicationRoute.Get("/application/:id/rental",
		auth_http.AuthorizedMiddleware(tokenMaker),
		CheckApplicationVisibilty(a.aService),
		a.getRentalByApplicationId(),
	)
}

func NewAdapter(lService listing_service.Service, aService application_service.Service) Adapter {
	return &adapter{
		lService: lService,
		aService: aService,
	}
}

func (a *adapter) preCreateApplications() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.PreCreateApplication
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		err := a.aService.PreCreateApplication(&payload, tkPayload.UserID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(payload)
	}
}

func (a *adapter) createApplications() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		var payload dto.CreateApplication
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}
		tkPayload, ok := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			payload.CreatorID = tkPayload.UserID
		} else {
			// validate key
			if payload.ApplicationKey == "" {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "application key is required"})
			}
			res, err := a.lService.VerifyApplicationLink(
				&listing_dto.VerifyApplicationLink{
					CreateApplicationLink: listing_dto.CreateApplicationLink{
						FullName:  payload.FullName,
						Email:     payload.Email,
						Phone:     payload.Phone,
						ListingId: payload.ListingID,
					},
					Key: payload.ApplicationKey,
				},
			)
			if err != nil || !res {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid application key"})
			}
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		// log.Println("CreateApplication payload:", payload)
		newApplication, err := a.aService.CreateApplication(&payload)
		if err != nil {
			if dbErr, ok := err.(*database.TXError); ok {
				return responses.DBTXErrorResponse(ctx, dbErr)
			}
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			if errors.Is(err, application_service.ErrAlreadyApplied) ||
				errors.Is(err, application_service.ErrListingIsClosed) ||
				errors.Is(err, application_service.ErrInvalidApplicant) {
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
		validator.RegisterValidation(dto.ApplicationFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, validation.GetValidationError(errs))
		}

		applications, err := a.aService.GetApplicationsByUserId(tkPayload.UserID, query)
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
		validator.RegisterValidation(dto.ApplicationFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		applications, err := a.aService.GetApplicationsToUser(tkPayload.UserID, query)
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

		application, err := a.aService.GetApplicationById(aid)
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
		validator.RegisterValidation(dto.ApplicationFieldsLocalKey, dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, *query); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, validation.GetValidationError(errs))
		}

		var userId uuid.UUID = uuid.Nil
		tkPayload, ok := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userId = tkPayload.UserID
		}

		res, err := a.aService.GetApplicationByIds(query.IDs, query.Fields, userId)
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

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		err = a.aService.UpdateApplicationStatus(aid, tkPayload.UserID, payload)
		if err != nil {
			if err == database.ErrRecordNotFound {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "application not found"})
			}
			if errors.Is(err, application_service.ErrInvalidStatusTransition) {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error()})
			}
			if errors.Is(err, application_service.ErrUnauthorizedUpdate) {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
			}
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (a *adapter) createApplicationMsgGroup() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		aid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid application id"})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		res, err := a.aService.CreateApplicationMsgGroup(aid, tkPayload.UserID)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "application not found"})
			}
			if errors.Is(err, application_service.ErrAnonymousApplicant) {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error()})
			}
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getApplicationMsgGroup() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		aid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid application id"})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		res, err := a.aService.GetApplicationMsgGroup(aid, tkPayload.UserID)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "application not found"})
			}
			if errors.Is(err, application_service.ErrAnonymousApplicant) {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error()})
			}
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getRentalByApplicationId() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		aid := ctx.Locals(ApplicationIdLocalKey).(int64)

		res, err := a.aService.GetRentalByApplicationId(aid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "rental not found"})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}
