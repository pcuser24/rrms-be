package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreateReminder struct {
	CreatorID uuid.UUID `json:"creatorId" validate:"required"`
	Title     string    `json:"title" validate:"required"`
	StartAt   time.Time `json:"startAt" validate:"required"`
	EndAt     time.Time `json:"endAt" validate:"required"`
	Note      *string   `json:"note" validate:"omitempty"`
	Location  string    `json:"location" validate:"required"`
}

func (c *CreateReminder) ToCreateReminderDB() database.CreateReminderParams {
	return database.CreateReminderParams{
		CreatorID: c.CreatorID,
		Title:     c.Title,
		StartAt:   c.StartAt,
		EndAt:     c.EndAt,
		Note:      types.StrN(c.Note),
		Location:  c.Location,
	}
}
