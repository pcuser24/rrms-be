package http

import (
	"github.com/gofiber/fiber/v2"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/statistic/service"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service service.Service
}

func NewAdapter(service service.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker) {
	statisticRoute := (*route).Group("/statistics")
	statisticRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))

	statisticRoute.Get("/properties", a.getPropertiesStatistic())
	statisticRoute.Get("/listings", a.getListingsStatistic())
	statisticRoute.Get("/applications", a.getApplicationStatistic())
	statisticRoute.Get("/payments", a.getPaymentsStatistic())
	statisticRoute.Get("/rentals", a.getRentalStatistic())
	statisticRoute.Get("/rentals/payments/arrears", a.getRentalPaymentArrearsStatistic())
	statisticRoute.Get("/rentals/payments/incomes", a.getRentalPaymentIncomesStatistic())

	landingRoute := (*route).Group("/landing")
	landingRoute.Get("/recent", a.getRecentListings())
	landingRoute.Get("/suggest", a.getListingSuggestions())
}
