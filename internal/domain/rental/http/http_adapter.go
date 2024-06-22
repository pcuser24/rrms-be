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
	rentalRoute.Post("/create/_pre", a.preCreateRental())
	rentalRoute.Post("/create/pre", a.createPreRental())
	// rentalRoute.Post("/create", a.createRental())
	rentalRoute.Get("/managed-rentals", a.getManagedRentals())
	rentalRoute.Get("/my-rentals", a.getMyRentals())
	rentalRoute.Group("/rental/:id").Use(GetRentalID(), CheckRentalVisibility(a.service))
	rentalRoute.Get("/rental/:id", a.getRental())
	rentalRoute.Patch("/rental/:id", a.updateRental())
	rentalRoute.Get("/rental/:id/contract", a.getRentalContract())
	rentalRoute.Get("/rental/:id/ping-contract", a.pingContract())
	rentalRoute.Post("/rental/:id/contract", a.createRentalContract())

	prerentalRoute := (*route).Group("/prerentals")
	prerentalRoute.Get("/to-me", auth_http.AuthorizedMiddleware(tokenMaker), a.getPreRentalsToMe())
	prerentalRoute.Get("/managed", auth_http.AuthorizedMiddleware(tokenMaker), a.getManagedPreRentals())
	prerentalRoute.Get("/prerental/:id",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		CheckPreRentalAccess(a.service),
		a.getPreRental(),
	)
	prerentalRoute.Patch("/prerental/:id/state",
		auth_http.GetAuthorizationMiddleware(tokenMaker),
		CheckPreRentalAccess(a.service),
		a.updatePreRentalState(),
	)

	contractRoute := (*route).Group("/contracts")
	contractRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	contractRoute.Get("/", a.getRentalContractsOfUser())
	contractRoute.Group("/contract/:id").Use(GetContractID())
	contractRoute.Get("/contract/:id", a.getRentalContract())
	// contractRoute.Patch("/contract/:id", a.updateContract())
	contractRoute.Patch("/contract/:id", a.updateContract())
	contractRoute.Patch("/contract/:id/content", a.updateContractContent())

	rentalPaymentRoute := (*route).Group("/rental-payments")
	rentalPaymentRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	rentalPaymentRoute.Get("/managed-payments", a.getManagedRentalPayments())
	rentalPaymentRoute.Post("/", a.createRentalPayment())
	rentalPaymentRoute.Get("/rental/:id",
		GetRentalID(),
		CheckRentalVisibility(a.service),
		a.getPaymentsOfRental(),
	)
	rentalPaymentRoute.Group("/rental-payment/:id").Use(GetRentalPaymentID())
	rentalPaymentRoute.Get("/rental-payment/:id", a.getRentalPayment())
	rentalPaymentRoute.Patch("/rental-payment/:id/plan", a.updatePlanRentalPayment())
	rentalPaymentRoute.Patch("/rental-payment/:id/issued", a.updateIssuedRentalPayment())
	rentalPaymentRoute.Patch("/rental-payment/:id/pending", a.updatePendingRentalPayment())
	rentalPaymentRoute.Patch("/rental-payment/:id/partiallypaid", a.updatePartiallyPaidRentalPayment())
	rentalPaymentRoute.Patch("/rental-payment/:id/payfine", a.updatePayfineRentalPayment())

	rentalComplaintRoute := (*route).Group("/rental-complaints")
	rentalComplaintRoute.Use(auth_http.AuthorizedMiddleware(tokenMaker))
	rentalComplaintRoute.Get("/", a.getRentalComplaintsOfUser())
	rentalComplaintRoute.Post("/create/_pre", a.preCreateRentalComplaint())
	rentalComplaintRoute.Post("/create", a.createRentalComplaint())
	rentalComplaintRoute.Group("/rental-complaint/:id").Use(GetRentalComplaintID())
	rentalComplaintRoute.Get("/rental-complaint/:id", a.getRentalComplaint())
	rentalComplaintRoute.Patch("/rental-complaint/:id", a.updateRentalComplaintStatus())
	rentalComplaintRoute.Get("/rental/:id/", a.getRentalComplaintsByRentalId())
	rentalComplaintRoute.Post("/rental-complaint/:id/replies/create/_pre", a.preCreateRentalComplaintReply())
	rentalComplaintRoute.Post("/rental-complaint/:id/replies/create", a.createRentalComplaintReply())
	rentalComplaintRoute.Get("/rental-complaint/:id/replies", a.getRentalComplaintReplies())
}
