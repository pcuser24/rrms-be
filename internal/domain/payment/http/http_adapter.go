package http

import (
	"github.com/gofiber/fiber/v2"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/payment/service"
	"github.com/user2410/rrms-backend/internal/domain/payment/service/vnpay"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	paymentService service.Service
}

func NewAdapter(paymentService service.Service) Adapter {
	return &adapter{
		paymentService: paymentService,
	}
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker) {
	paymentRoute := (*route).Group("/payments")
	// paymentRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	paymentRoute.Get("/my-payments", auth_http.AuthorizedMiddleware(tokenMaker), a.getMyPayments())
	paymentRoute.Get("/payment/:id", auth_http.AuthorizedMiddleware(tokenMaker), a.getPaymentById())

	_, ok := a.paymentService.(*vnpay.VnPayService)
	if ok {
		vnpayRoute := paymentRoute.Group("/vnpay")
		vnpayRoute.Post("/create_payment_url/:paymentId", auth_http.AuthorizedMiddleware(tokenMaker), a.vnpCreatePaymentUrl())
		vnpayRoute.Get("/vnpay_return", a.vnpReturn())
		vnpayRoute.Get("/vnpay_ipn", a.vnpIpn())
		vnpayRoute.Post("/querydr", a.vnpQuerydr())
		vnpayRoute.Post("/refund", a.vnpRefund())
	}

}
