package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/reminder"
	"github.com/user2410/rrms-backend/internal/domain/reminder/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service reminder.Service
}

func NewAdapter(service reminder.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker) {
	reminderRoute := (*route).Group("/reminders")
	reminderRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	reminderRoute.Get("/", a.getRemindersOfUser())
	reminderRoute.Group("/reminder/:id").Use(GetReminderId())
	reminderRoute.Get("/reminder/:id",
		CheckReminderVisibility(a.service),
		a.getReminder(),
	)
	reminderRoute.Patch("/reminder/:id",
		CheckReminderVisibility(a.service),
		a.updateReminderStatus(),
	)
}

func (a *adapter) getReminder() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		reminder, err := a.service.GetReminderById(id)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(reminder)
	}
}

func (a *adapter) getRemindersOfUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var query dto.GetRemindersQuery
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		res, err := a.service.GetRemindersOfUser(tkPayload.UserID, &query)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) updateReminderStatus() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		status := ctx.Query("status")
		if status == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "status is required"})
		}

		err = a.service.UpdateReminderStatus(id, database.REMINDERSTATUS(status))
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
