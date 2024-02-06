package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/user2410/rrms-backend/internal/domain/auth"

	"github.com/user2410/rrms-backend/internal/utils/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service auth.Service
}

func NewAdapter(service auth.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	authRoute := (*router).Group("/auth")

	credentialGroup := authRoute.Group("/credential")
	credentialGroup.Post("/register", a.credentialRegister())
	credentialGroup.Post("/login", GetAuthorizationMiddleware(tokenMaker), a.credentialLogin())
	credentialGroup.Put("/refresh", a.credentialRefresh())
	credentialGroup.Delete("/logout", AuthorizedMiddleware(tokenMaker), a.credentialLogout())

	bffGroup := authRoute.Group("/bff")
	bffUserGroup := bffGroup.Group("/user")
	bffUserGroup.Post("/create", a.bffCreateUser())
	bffUserGroup.Get("/get-by-id/:id", a.bffGetUserById())
	bffUserGroup.Get("/get-by-email/:email", a.bffGetUserByEmail())
	bffUserGroup.Get("/get-by-account/:id", a.bffGetUserByAccount())
	bffUserGroup.Patch("/update/:id", a.bffUpdateUser())
	bffUserGroup.Delete("/delete/:id", a.bffDeleteUser())
	bffGroup.Patch("link-account", a.bffLinkAccount())
	bffGroup.Patch("unlink-account", a.bffUnlinkAccount())
	bffSessionGroup := bffGroup.Group("/session")
	bffSessionGroup.Post("/create", a.bffCreateSession())
	bffSessionGroup.Get("/user/:token", a.bffGetSessionAndUser())
	bffSessionGroup.Patch("/update", a.bffUpdateSession())
	bffSessionGroup.Delete("/delete/:token", a.bffDeleteSession())
	bffGroup.Post("/create-verification-token", a.bffCreateVerificationToken())
	bffGroup.Put("/use-verification-token", a.bffUseVerificationToken())
}

func (a *adapter) credentialRegister() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.RegisterUser
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
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

			if errors.Is(err, auth.ErrInvalidCredential) {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
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
			case auth.ErrInvalidCredential, auth.ErrInvalidSession:
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

func (a *adapter) bffCreateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffGetUserByEmail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Params("email")

		user, err := a.service.GetUserByEmail(email)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("User with email %s not found", email)})
			}
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.JSON(user)
	}
}

func (a *adapter) bffGetUserById() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid id"})
		}

		user, err := a.service.GetUserById(uid)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("User with id %s not found", id)})
			}
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.JSON(user)
	}
}

func (a *adapter) bffGetUserByAccount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffUpdateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffDeleteUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffLinkAccount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffUnlinkAccount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffCreateSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffGetSessionAndUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffUpdateSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffDeleteSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffCreateVerificationToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *adapter) bffUseVerificationToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
