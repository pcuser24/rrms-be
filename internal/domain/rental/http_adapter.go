package rental

import (
	"github.com/gofiber/fiber/v2"
)

type Adapter interface {
	RegisterServer(route *fiber.Router)
}

type adapter struct {
	service Service
}

func NewAdapter(service Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(route *fiber.Router) {
	rentalRoute := (*route).Group("/rental")
	rentalRoute.Get("/policies", a.getRentalPolicies())
}

func (a *adapter) getRentalPolicies() fiber.Handler {
	return func(c *fiber.Ctx) error {
		res, err := a.service.GetAllRentalPolicies()
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"items": res,
		})
	}
}
