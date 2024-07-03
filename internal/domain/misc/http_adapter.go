package misc

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/misc/dto"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/types"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service misc_service.Service
}

func NewAdapter(service misc_service.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	miscRoute := (*router).Group("/misc")

	notificationRoute := miscRoute.Group("/notifications").Use(auth_http.AuthorizedMiddleware(tokenMaker))

	notificationRoute.Get("/", a.getNotificationsOfUser())
	// notificationRoute.Get("/notification/:id", a.getNotification())
	notificationRoute.Patch("/notification/:id/seend", a.updateNotificationStatus())

	deviceRoute := notificationRoute.Group("/devices")
	deviceRoute.Post("/", a.createNotificationDevice())
	deviceRoute.Get("/", a.getNotificationDevices())

}

func (a *adapter) createNotificationDevice() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.CreateNotificationDevice
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		res, err := a.service.CreateNotificationDevice(tkPayload.UserID, &payload)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) getNotificationDevices() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sessionId, err := uuid.Parse(ctx.Query("sessionId"))
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid session id"})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)

		res, err := a.service.GetNotificationDevice(tkPayload.UserID, sessionId, ctx.Query("token"), ctx.Query("platform"))
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.SendStatus(fiber.StatusNotFound)
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getNotificationsOfUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var q dto.GetNotificationsOfUserQuery
		if err := ctx.QueryParser(&q); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		res, err := a.service.GetNotificationsOfUser(tkPayload.UserID, q)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

// func (a *adapter) getNotification() fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		return nil
// 	}
// }

func (a *adapter) updateNotificationStatus() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid notification id"})
		}

		// tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		// TODO: check user permission

		params := dto.UpdateNotification{
			ID:   id,
			Seen: types.Ptr(true),
		}

		err = a.service.UpdateNotification(&params)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
