package server

import (
	"log"

	application_http "github.com/user2410/rrms-backend/internal/domain/application/http"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/chat"
	listing_http "github.com/user2410/rrms-backend/internal/domain/listing/http"
	"github.com/user2410/rrms-backend/internal/domain/misc"
	payment_http "github.com/user2410/rrms-backend/internal/domain/payment/http"
	property_http "github.com/user2410/rrms-backend/internal/domain/property/http"
	reminder_http "github.com/user2410/rrms-backend/internal/domain/reminder/http"
	rental_http "github.com/user2410/rrms-backend/internal/domain/rental/http"
	statistic_http "github.com/user2410/rrms-backend/internal/domain/statistic/http"
	unit_http "github.com/user2410/rrms-backend/internal/domain/unit/http"
)

func (c *serverCommand) setupHttpServer() {
	apiRoute := c.httpServer.GetApiRoute()
	// v1 := (*apiRoute).Group("/v1")apiRoute

	auth_http.
		NewAdapter(c.internalServices.AuthService).
		RegisterServer(apiRoute, c.tokenMaker)
	property_http.
		NewAdapter(c.internalServices.PropertyService).
		RegisterServer(apiRoute, c.tokenMaker)
	unit_http.NewAdapter(c.internalServices.UnitService, c.internalServices.PropertyService).
		RegisterServer(apiRoute, c.tokenMaker)
	listing_http.
		NewAdapter(c.internalServices.ListingService, c.internalServices.PropertyService, c.internalServices.UnitService).
		RegisterServer(apiRoute, c.tokenMaker)
	rental_http.
		NewAdapter(c.internalServices.RentalService).
		RegisterServer(apiRoute, c.tokenMaker)
	application_http.
		NewAdapter(c.internalServices.ListingService, c.internalServices.ApplicationService).
		RegisterServer(apiRoute, c.tokenMaker)
	payment_http.
		NewAdapter(c.internalServices.PaymentService).
		RegisterServer(apiRoute, c.tokenMaker)
	chat.
		NewWSChatAdapter(c.internalServices.ChatService).
		RegisterServer(c.httpServer.GetFibApp(), c.tokenMaker)
	chat.
		NewHttpAdapter(c.internalServices.ChatService).
		RegisterServer(apiRoute, c.tokenMaker)
	reminder_http.
		NewAdapter(c.internalServices.ReminderService).
		RegisterServer(apiRoute, c.tokenMaker)
	statistic_http.
		NewAdapter(c.internalServices.StatisticService).
		RegisterServer(apiRoute, c.tokenMaker)
	misc.
		NewAdapter(c.internalServices.MiscService).
		RegisterServer(apiRoute, c.tokenMaker)
}

func (c *serverCommand) runHttpServer(errChan chan error) {
	log.Println("Starting HTTP server...")
	var port uint16 = 8080
	if c.config.Port != nil {
		port = *c.config.Port
	}
	if err := c.httpServer.Start(port); err != nil {
		errChan <- err
	}
}
