package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type CreateReminder struct {
	Title          string                          `json:"title" validate:"required"`
	StartAt        time.Time                       `json:"startAt" validate:"required"`
	EndAt          time.Time                       `json:"endAt" validate:"required"`
	Note           *string                         `json:"note" validate:"required"`
	Location       string                          `json:"location" validate:"required"`
	Priority       *int32                          `json:"priority" validate:"omitempty"`
	RecurrenceMode database.REMINDERRECURRENCEMODE `json:"recurrenceMode" validate:"required"`
	Members        []uuid.UUID
}
