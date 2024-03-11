package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/payment/service"
	"github.com/user2410/rrms-backend/internal/domain/payment/service/vnpay"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	paymentService service.Service
	vnpayService   *vnpay.Service
}

func NewAdapter(paymentService service.Service, vnpayService *vnpay.Service) Adapter {
	return &adapter{
		paymentService: paymentService,
		vnpayService:   vnpayService,
	}
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker) {
	paymentRoute := (*route).Group("/payments")
	// paymentRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	paymentRoute.Get("/payment/:id",
		auth_http.AuthorizedMiddleware(tokenMaker),
		a.getPaymentById(),
	)

	vnpayRoute := paymentRoute.Group("/vnpay")
	vnpayRoute.Post("/create_payment_url/:paymentId", auth_http.AuthorizedMiddleware(tokenMaker), a.vnpCreatePaymentUrl())
	vnpayRoute.Get("/vnpay_return", a.vnpReturn())
	vnpayRoute.Get("/vnpay_ipn", a.vnpIpn())
	vnpayRoute.Post("/querydr", a.vnpQuerydr())
	vnpayRoute.Post("/refund", a.vnpRefund())
}

func (a *adapter) getPaymentById() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid payment id",
			})
		}

		payment, err := a.paymentService.GetPaymentById(id)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Payment not found",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		return c.Status(fiber.StatusOK).JSON(payment)
	}
}
