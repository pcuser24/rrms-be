package unit

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service Service
}

func NewAdapter(service Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(router fiber.Router, tokenMaker token.Maker) {
	unitRoute := router.Group("/unit")

	unitRoute.Use(auth.NewAuthMiddleware(tokenMaker))

	unitRoute.Post("/create", a.createUnit())
	unitRoute.Get("/get-by-id/:id", a.getUnitById())
}

func (a *adapter) createUnit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.SendString("create unit")
		return nil
	}
}

func (a *adapter) getUnitById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
