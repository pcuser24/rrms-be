package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

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
