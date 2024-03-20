package server

import (
	"log"

	application_http "github.com/user2410/rrms-backend/internal/domain/application/http"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/chat"
	listing_http "github.com/user2410/rrms-backend/internal/domain/listing/http"
	payment_http "github.com/user2410/rrms-backend/internal/domain/payment/http"
	property_http "github.com/user2410/rrms-backend/internal/domain/property/http"
	"github.com/user2410/rrms-backend/internal/domain/rental"
	"github.com/user2410/rrms-backend/internal/domain/storage"
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
	rental.
		NewAdapter(c.internalServices.RentalService).
		RegisterServer(apiRoute)
	application_http.
		NewAdapter(c.internalServices.ListingService, c.internalServices.ApplicationService).
		RegisterServer(apiRoute, c.tokenMaker)
	storage.
		NewAdapter(c.internalServices.StorageService).
		RegisterServer(apiRoute, c.tokenMaker)
	payment_http.
		NewAdapter(c.internalServices.PaymentService, c.internalServices.VnpService).
		RegisterServer(apiRoute, c.tokenMaker)
	chat.
		NewWSChatAdapter(c.internalServices.ChatService).
		RegisterServer(c.httpServer.GetFibApp(), c.tokenMaker)
	chat.
		NewHttpAdapter(c.internalServices.ChatService).
		RegisterServer(apiRoute, c.tokenMaker)
}

func (c *serverCommand) runHttpServer(errChan chan error) {
	log.Println("Starting HTTP server...")
	var port uint16 = 8000
	if c.config.Port != nil {
		port = *c.config.Port
	}
	if err := c.httpServer.Start(port); err != nil {
		errChan <- err
	}
}
