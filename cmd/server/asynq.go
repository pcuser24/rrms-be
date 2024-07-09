package server

import (
	"log"

	application_asynctask "github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	property_asynctask "github.com/user2410/rrms-backend/internal/domain/property/asynctask"
	rental_asynctask "github.com/user2410/rrms-backend/internal/domain/rental/asynctask"
)

func (c *serverCommand) setupAsyncTaskProcessor() {
	property_asynctask.
		NewAdapter(c.internalServices.PropertyService).
		Register(c.asyncTaskProcessor)
	application_asynctask.
		NewAdapter(c.internalServices.ApplicationService).
		Register(c.asyncTaskProcessor)
	rental_asynctask.
		NewAdapter(c.internalServices.RentalService).
		Register(c.asyncTaskProcessor)
}

func (c *serverCommand) runAsyncTaskProcessor(errChan chan error) {
	log.Println("Starting async task processor...")
	if err := c.asyncTaskProcessor.Start(); err != nil {
		errChan <- err
	}
}
