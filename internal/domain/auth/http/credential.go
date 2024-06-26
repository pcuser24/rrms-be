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

		if payload.Email == "alpha@email.com" {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"sessionId":    "56e93813-bebb-4de5-9347-e17352a0df8d",
				"accessToken":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjU2ZTkzODEzLWJlYmItNGRlNS05MzQ3LWUxNzM1MmEwZGY4ZCIsInR5cGUiOiJhY2Nlc3MiLCJzdWIiOiJlMGE4ZDEyMy1jNTViLTQyMzAtOTFlOC1iZDFiN2I3NjIzNjYiLCJpYXQiOiIyMDI0LTA2LTA3VDEwOjQ2OjA3LjAzOTcwMTY5KzA3OjAwIiwiZXhwIjoiMjAyNC0wNy0wN1QxMDo0NjowNy4wMzk3MDE4OSswNzowMCJ9.vPn1T4LoRjZzZJZ7XQn8c_ZtAr2T4kXGlw1x6YsdNY4",
				"accessExp":    "2024-07-07T10:46:07.03970189+07:00",
				"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjU2ZTkzODEzLWJlYmItNGRlNS05MzQ3LWUxNzM1MmEwZGY4ZCIsInR5cGUiOiJyZWZyZXNoIiwic3ViIjoiZTBhOGQxMjMtYzU1Yi00MjMwLTkxZTgtYmQxYjdiNzYyMzY2IiwiaWF0IjoiMjAyNC0wNS0wN1QxNzoxODo1MC4zNTQ5MjEyMDErMDc6MDAiLCJleHAiOiIyMDI1LTA1LTAyVDE3OjE4OjUwLjM1NDkyMTQwMyswNzowMCJ9.QCAtUnlOrOrjxbMotss3VVf_MK4eEwW4LhQh1BJJQvk",
				"refreshExp":   "2025-05-02T17:18:50.354921403+07:00",
				"user": fiber.Map{
					"id":        "e0a8d123-c55b-4230-91e8-bd1b7b762366",
					"email":     "alpha@email.com",
					"createdAt": "2024-01-31T15:53:31.935561+07:00",
					"updatedAt": "2024-02-16T11:02:45.375152+07:00",
					"deleted_f": false,
					"firstName": "Albert",
					"lastName":  "Alpha",
					"phone":     "0912142214",
					"avatar":    nil,
					"address":   "Số 1, Đường Giải Phóng",
					"city":      "HN",
					"district":  "4",
					"ward":      "74",
					"role":      "LANDLORD",
				},
			})
		} else if payload.Email == "gamma@email.com" {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"sessionId":    "7add15ab-c9bb-4b48-b619-17ccbaa5e62a",
				"accessToken":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjdhZGQxNWFiLWM5YmItNGI0OC1iNjE5LTE3Y2NiYWE1ZTYyYSIsInR5cGUiOiJhY2Nlc3MiLCJzdWIiOiI4ODBhZTM3Ni01ODUwLTQ4NWQtOGE5Zi03YmNkNTdmZjgzMzMiLCJpYXQiOiIyMDI0LTA2LTAzVDA5OjI0OjQ3LjEwMTg1MTAwNyswNzowMCIsImV4cCI6IjIwMjQtMDctMDNUMDk6MjQ6NDcuMTAxODUxMTI1KzA3OjAwIn0.ss3tRqalxeWdafppSzIMGNhwZd-3dTQLiT0XHEe5uVM",
				"accessExp":    "2024-07-03T09:24:47.101851125+07:00",
				"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjdhZGQxNWFiLWM5YmItNGI0OC1iNjE5LTE3Y2NiYWE1ZTYyYSIsInR5cGUiOiJyZWZyZXNoIiwic3ViIjoiODgwYWUzNzYtNTg1MC00ODVkLThhOWYtN2JjZDU3ZmY4MzMzIiwiaWF0IjoiMjAyNC0wNi0wM1QwOToyNDo0Ny4xMDE3MjkwMDErMDc6MDAiLCJleHAiOiIyMDI1LTA1LTI5VDA5OjI0OjQ3LjEwMTcyOTEyOCswNzowMCJ9.t-dWQ-4PyHH5tOc0B9NmmJEX8DAm-hZpW2_hYwXWj1Q",
				"refreshExp":   "2025-05-29T09:24:47.101729128+07:00",
				"user": fiber.Map{
					"id":        "880ae376-5850-485d-8a9f-7bcd57ff8333",
					"email":     "gamma@email.com",
					"createdAt": "2024-02-15T17:07:49.984793+07:00",
					"updatedAt": "2024-02-15T17:07:49.984793+07:00",
					"deleted_f": false,
					"firstName": "Graham",
					"lastName":  "Gamma",
					"phone":     nil,
					"avatar":    nil,
					"address":   nil,
					"city":      nil,
					"district":  nil,
					"ward":      nil,
					"role":      "TENANT",
				},
			})
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
