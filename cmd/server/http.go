package server

import (
	"log"
	"time"

	application_http "github.com/user2410/rrms-backend/internal/domain/application/http"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	listing_http "github.com/user2410/rrms-backend/internal/domain/listing/http"
	property_http "github.com/user2410/rrms-backend/internal/domain/property/http"
	"github.com/user2410/rrms-backend/internal/domain/rental"
	"github.com/user2410/rrms-backend/internal/domain/storage"
	unit_http "github.com/user2410/rrms-backend/internal/domain/unit/http"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (c *serverCommand) setupHttpServer() {
	c.httpServer = http.NewServer(
		fiber.Config{
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
		cors.Config{
			AllowOrigins: c.config.AllowOrigins,
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		},
	)
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
	listing_http.NewAdapter(c.internalServices.ListingService, c.internalServices.PropertyService, c.internalServices.UnitService).
		RegisterServer(apiRoute, c.tokenMaker)
	rental.
		NewAdapter(c.internalServices.RentalService).
		RegisterServer(apiRoute)
	application_http.
		NewAdapter(c.internalServices.ApplicationService).
		RegisterServer(apiRoute, c.tokenMaker)
	storage.
		NewAdapter(c.internalServices.StorageService).
		RegisterServer(apiRoute, c.tokenMaker)
}

func (c *serverCommand) runHttpServer() {
	log.Println("Starting HTTP server...")
	if err := c.httpServer.Start(8000); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
