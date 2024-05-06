package http

import (
	"github.com/gofiber/fiber/v2"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/rental/service"
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
	rentalRoute := (*route).Group("/rentals")
	rentalRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	rentalRoute.Post("/", a.createRental())
	rentalRoute.Get("/rental/:id", CheckRentalVisibility(a.service), a.getRental())
	rentalRoute.Patch("/rental/:id", CheckRentalVisibility(a.service), a.updateRental())
	rentalRoute.Get("/rental/:id/contract", CheckRentalVisibility(a.service), a.getRentalContract())
	rentalRoute.Get("/rental/:id/ping-contract", CheckRentalVisibility(a.service), a.pingContract())
	rentalRoute.Post("/rental/:id/contract", CheckRentalVisibility(a.service), a.createRentalContract())

	contractRoute := (*route).Group("/contracts")
	contractRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	contractRoute.Get("/contract/:id", a.getContract())
	// contractRoute.Patch("/contract/:id", a.updateContract())
	contractRoute.Patch("/contract/:id", a.updateContract())
	contractRoute.Patch("/contract/:id/content", a.updateContractContent())

	rentalPaymentRoute := (*route).Group("/rental-payments")
	rentalPaymentRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	rentalPaymentRoute.Post("/", a.createRentalPayment())
	rentalPaymentRoute.Get("/rental-payment/:id", a.getRentalPayment())
	rentalPaymentRoute.Get("/rental/:id", CheckRentalVisibility(a.service), a.getPaymentsOfRental())
	rentalPaymentRoute.Patch("/rental-payment/:id/plan", a.updatePlanRentalPayment())
	rentalPaymentRoute.Patch("/rental-payment/:id/issued", a.updateIssuedRentalPayment())
	rentalPaymentRoute.Patch("/rental-payment/:id/pending", a.updatePendingRentalPayment())

	rentalComplaintRoute := (*route).Group("/rental-complaints")
	rentalComplaintRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	rentalComplaintRoute.Post("/", a.createRentalComplaint())
	rentalComplaintRoute.Get("/rental-complaint/:id", a.getRentalComplaint())
	rentalComplaintRoute.Patch("/rental-complaint/:id", a.updateRentalComplaintStatus())
	rentalComplaintRoute.Get("/rental/:id/", a.getRentalComplaintsByRentalId())
	rentalComplaintRoute.Post("/rental-complaint/:id/replies", a.createRentalComplaintReply())
	rentalComplaintRoute.Get("/rental-complaint/:id/replies", a.getRentalComplaintReplies())
}
