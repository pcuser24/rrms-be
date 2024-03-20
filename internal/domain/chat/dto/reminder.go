package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type IncomingUpdateReminderEvent struct {
	ID     int64                   `json:"id" validate:"required"`
	Status database.REMINDERSTATUS `json:"status" validate:"required"`
}

type OutgoingUpdateReminderEvent struct {
	ID     int64                   `json:"id" validate:"required"`
	UserId uuid.UUID               `json:"userId" validate:"required"`
	Status database.REMINDERSTATUS `json:"status" validate:"required"`
}
