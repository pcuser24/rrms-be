package http

import (
	"github.com/gofiber/fiber/v2"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/statistic/service"
	unit_http "github.com/user2410/rrms-backend/internal/domain/unit/http"
	unit_service "github.com/user2410/rrms-backend/internal/domain/unit/service"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker, unitService unit_service.Service)
}

type adapter struct {
	service service.Service
}

func NewAdapter(service service.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker, unitService unit_service.Service) {
	statisticRoute := (*route).Group("/statistics")
	statisticRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))

	managerStatisticRoute := statisticRoute.Group("/manager")
	managerStatisticRoute.Get("/properties", a.getPropertiesStatistic())
	managerStatisticRoute.Get("/applications", a.getApplicationStatistic())
	managerStatisticRoute.Get("/payments", a.getPaymentsStatistic())
	managerStatisticRoute.Get("/maintenance", a.getManagerMaintenanceStatistic())
	managerStatisticRoute.Get("/tenants", a.getTotalTenantsStatistic())
	managerStatisticRoute.Get("/rentals/payments/arrears", a.getRentalPaymentArrearsStatistic())
	managerStatisticRoute.Get("/rentals/payments/incomes", a.getRentalPaymentIncomesStatistic())

	tenantStatisticRoute := statisticRoute.Group("/tenant")
	tenantStatisticRoute.Get("/rentals", a.getTenantRentalStatistic())
	tenantStatisticRoute.Get("/maintenances", a.getTenantMaintenanceStatistic())
	tenantStatisticRoute.Get("/expenditures", a.getTenantExpenditureStatistic())
	tenantStatisticRoute.Get("/arrears", a.getTenantArrearsStatistic())

	// public routes
	landingRoute := (*route).Group("/landing")
	landingRoute.Get("/recent", a.getRecentListings())
	landingRoute.Get("/suggest/listings/listing/:id", a.getListingSuggestions())
	landingRoute.Get("/rentals/unit/:id", auth_http.GetAuthorizationMiddleware(tokenMaker), unit_http.CheckUnitVisiblitiy(unitService), a.getTotalTenantsOfUnitStatistic())
}
