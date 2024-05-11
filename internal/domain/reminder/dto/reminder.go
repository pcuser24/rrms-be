package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreateReminder struct {
	CreatorID       uuid.UUID                       `json:"creatorId" validate:"required"`
	Title           string                          `json:"title" validate:"required"`
	StartAt         time.Time                       `json:"startAt" validate:"required"`
	EndAt           time.Time                       `json:"endAt" validate:"required"`
	Note            *string                         `json:"note" validate:"omitempty"`
	Location        string                          `json:"location" validate:"required"`
	Priority        *int32                          `json:"priority" validate:"omitempty"`
	RecurrenceDay   *int32                          `json:"recurrenceDay" validate:"omitempty,min=1,max=28"`
	RecurrenceMonth *int32                          `json:"recurrenceMonth" validate:"omitempty,min=1,max=12"`
	RecurrenceMode  database.REMINDERRECURRENCEMODE `json:"recurrenceMode" validate:"required,oneof=NONE WEEKLY MONTHLY"`
	ResourceTag     string                          `json:"resourceTag" validate:"required"`
	Members         []uuid.UUID                     `json:"members" validate:"required,dive"`
}

func (c *CreateReminder) ToCreateReminderDB() database.CreateReminderParams {
	return database.CreateReminderParams{
		CreatorID:       c.CreatorID,
		Title:           c.Title,
		StartAt:         c.StartAt,
		EndAt:           c.EndAt,
		Note:            types.StrN(c.Note),
		Location:        c.Location,
		Priority:        types.Int32N(c.Priority),
		RecurrenceDay:   types.Int32N(c.RecurrenceDay),
		RecurrenceMonth: types.Int32N(c.RecurrenceMonth),
		RecurrenceMode:  c.RecurrenceMode,
		ResourceTag:     c.ResourceTag,
	}
}
