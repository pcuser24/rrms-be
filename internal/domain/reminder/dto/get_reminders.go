package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type GetRemindersQuery struct {
	CreatorID       uuid.UUID                       `query:"creatorId" validate:"omitempty"`
	MinStartAt      time.Time                       `query:"minStartAt" validate:"omitempty"`
	MaxStartAt      time.Time                       `query:"maxStartAt" validate:"omitempty"`
	MinEndAt        time.Time                       `query:"minEndAt" validate:"omitempty"`
	MaxEndAt        time.Time                       `query:"maxEndAt" validate:"omitempty"`
	Priority        *int32                          `query:"priority" validate:"omitempty"`
	Status          database.REMINDERSTATUS         `query:"status" validate:"omitempty"`
	RecurrenceMode  database.REMINDERRECURRENCEMODE `query:"recurrenceMode" validate:"omitempty"`
	RecurrenceDay   *int32                          `query:"recurrenceDay" validate:"omitempty"`
	RecurrenceMonth *int32                          `query:"recurrenceMonth" validate:"omitempty"`
	ResourceTag     *string                         `query:"resourceTag" validate:"omitempty"`

	Members []uuid.UUID `query:"members" validate:"omitempty,dive"`
}
