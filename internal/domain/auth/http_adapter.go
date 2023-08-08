package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils"
)

type Adapter interface {
	RegisterServer(route *fiber.Router)
}

type adapter struct {
	service AuthService
}

func NewAdapter(service AuthService) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router) {
	authRoute := (*router).Group("/auth")

	credentialGroup := authRoute.Group("/credential")
	credentialGroup.Post("/register", a.credentialRegisterHandle())
	credentialGroup.Post("/login", a.credentialLoginHandle())

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

func (a *adapter) credentialRegisterHandle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.RegisterUser
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetFibValidationError(errs)})
		}

		res, err := a.service.RegisterUser(&payload)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				responses.DBErrorResponse(ctx, pgErr)
				return nil
			}

			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			return nil
		}

		return ctx.JSON(fiber.Map{
			"accessToken": res.AccessToken,
			"accessExp":   res.AccessPayload.ExpiredAt,
			"user":        res.User.ToUserResponse(),
		})
	}
}

func (a *adapter) credentialLoginHandle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.LoginUser
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if errs := utils.ValidateStruct(payload); len(errs) > 0 && errs[0].Error {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.GetFibValidationError(errs)})
		}

		res, err := a.service.Login(&payload)
		if err != nil {
			// check if error is database error
			if dbErr, ok := err.(*pgconn.PgError); ok {
				responses.DBErrorResponse(ctx, dbErr)
				return nil
			}

			if errors.Is(err, sql.ErrNoRows) {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
			}

			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.JSON(fiber.Map{
			"accessToken": res.AccessToken,
			"accessExp":   res.AccessPayload.ExpiredAt,
			"user":        res.User.ToUserResponse(),
		})
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
			c.Status(http.StatusInternalServerError)
			fmt.Println(err)
			return err
		}

		return c.JSON(user)
	}
}

func (a *adapter) bffGetUserById() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid id")
		}

		user, err := a.service.GetUserById(uid)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			fmt.Println(err)
			return err
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
