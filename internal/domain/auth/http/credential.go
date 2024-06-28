package http

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	auth_service "github.com/user2410/rrms-backend/internal/domain/auth/service"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"

	"github.com/user2410/rrms-backend/internal/utils/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

func (a *adapter) credentialRegister() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.RegisterUser
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.Register(&payload)
		if err != nil {
			dbErrCode := database.ErrorCode(err)
			if dbErrCode == database.UniqueViolation {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "email already exists"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func (a *adapter) credentialLogin() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		admin := ctx.Query("admin")

		var payload dto.LoginUser
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		session := &dto.CreateSession{
			ID:        uuid.Nil,
			UserAgent: ctx.Context().UserAgent(),
			ClientIp:  ctx.IP(),
		}
		tkPayload, ok := ctx.Locals(AuthorizationPayloadKey).(*token.Payload)
		if ok {
			session.ID = tkPayload.ID
		}

		res, err := a.service.Login(&payload, session)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"message": "no user with such email"})
			}

			if errors.Is(err, auth_service.ErrInvalidCredential) {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
			}

			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		if admin != "" && res.User.Role != database.USERROLEADMIN {
			return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"message": "forbidden"})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) credentialGetCurrentUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenPayload := ctx.Locals(AuthorizationPayloadKey).(*token.Payload)

		user, err := a.service.GetUserById(tokenPayload.UserID)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
			}
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(user)
	}

}

func (a *adapter) credentialGetUserByIds() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var q struct {
			Ids []uuid.UUID `json:"ids" validate:"required,dive,uuid4"`
		}
		if err := ctx.QueryParser(&q); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		res, err := a.service.GetUserByIds(q.Ids)
		if err != nil {
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) credentialRefresh() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.RefreshToken
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		res, err := a.service.RefreshAccessToken(payload.AccessToken, payload.RefreshToken)
		if err != nil {
			switch err {
			case auth_service.ErrInvalidCredential, auth_service.ErrInvalidSession:
				return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
			case token.ErrInvalidToken:
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
			default:
				return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			}
		}

		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"accessToken": res.AccessToken,
			"accessExp":   res.AccessExp,
		})
	}
}

func (a *adapter) credentialUpdateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenPayload := ctx.Locals(AuthorizationPayloadKey).(*token.Payload)

		var payload dto.UpdateUser
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		err := a.service.UpdateUser(tokenPayload.UserID, tokenPayload.UserID, &payload)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (a *adapter) credentialLogout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenPayload := c.Locals(AuthorizationPayloadKey).(*token.Payload)
		err := a.service.Logout(tokenPayload.ID)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return nil
	}
}
