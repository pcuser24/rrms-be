package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type SendEmailOnNewApplicationPayload struct {
	Email         string    `json:"email" validate:"required,email"`
	Username      string    `json:"username"`
	ApplicationId int64     `json:"applicationId"`
	ListingId     uuid.UUID `json:"listingId"`
}

type UpdateApplicationStatusPayload struct {
	Email         string                     `json:"email" validate:"required,email"`
	ApplicationId int64                      `json:"applicationId"`
	NewStatus     database.APPLICATIONSTATUS `json:"newStatus"`
	OldStatus     database.APPLICATIONSTATUS `json:"oldStatus"`
	Message       *string                    `json:"message"`
}
