package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/chat/dto"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

type HttpAdapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type httpAdapter struct {
	service Service
}

func NewHttpAdapter(s Service) HttpAdapter {
	return &httpAdapter{
		service: s,
	}
}

func (a *httpAdapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	chatRoute := (*router).Group("/chat")

	chatRoute.Get("/group/:id/messages",
		auth_http.AuthorizedMiddleware(tokenMaker),
		CheckGroupMembership(a.service),
		a.GetMessagesOfGroup(),
	)
}

func (a *httpAdapter) GetMessagesOfGroup() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var query dto.GetMessagesOfGroupQuery
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		groupId := ctx.Locals(GroupIDLocalKey).(int64)
		messages, err := a.service.GetMessagesOfGroup(groupId, query.Offset, query.Limit)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(messages)
	}
}
