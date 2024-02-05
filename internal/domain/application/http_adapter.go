package application

import (
	"log"
	"strconv"

	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service Service
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker) {
	applicationRoute := (*route).Group("/applications")

	applicationRoute.Use(auth.AuthorizedMiddleware(tokenMaker))

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
	applicationRoute.Put("/application/:id",
		// TODO: A middleware to check if the user is a property manager
		a.updateApplicationStatus(),
	)
	applicationRoute.Delete("/application/:id", a.deleteApplication())

}

func NewAdapter(service Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) createApplications() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		var payload dto.CreateApplication
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}
		payload.CreatorID = tkPayload.UserID
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		log.Println("CreateApplication payload:", payload)
		newApplication, err := a.service.CreateApplication(&payload)
		if err != nil {
			if dbErr, ok := err.(*database.TXError); ok {
				return responses.DBTXErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(newApplication)
	}
}

func (a *adapter) getMyApplications() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		applications, err := a.service.GetApplicationsByUserId(tkPayload.UserID)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(applications)
	}
}

// Get applications to properties that I manage
func (a *adapter) getApplicationsToMe() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		applications, err := a.service.GetApplicationsToUser(tkPayload.UserID)
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

func (a *adapter) updateApplicationStatus() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		aid, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid application id"})
		}

		status := ctx.Query("status")
		if status == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "status is required"})
		}

		err = a.service.UpdateApplicationStatus(aid, database.APPLICATIONSTATUS(status))
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
