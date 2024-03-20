package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type CreateReminder struct {
	Title    string    `json:"title" validate:"required"`
	StartAt  time.Time `json:"startAt" validate:"required"`
	EndAt    time.Time `json:"endAt" validate:"required"`
	Note     *string   `json:"note" validate:"required"`
	Location string    `json:"location" validate:"required"`
	Members  []uuid.UUID
}

type UpdateReminderStatus struct {
	ID     int64                   `json:"id" validate:"required"`
	Status database.REMINDERSTATUS `json:"status" validate:"required"`
}
