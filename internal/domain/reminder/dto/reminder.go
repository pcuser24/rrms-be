package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

type UpdateReminder struct {
	ID              int64                           `json:"id"`
	Title           *string                         `json:"title"`
	StartAt         time.Time                       `json:"startAt"`
	EndAt           time.Time                       `json:"endAt"`
	Note            *string                         `json:"note"`
	Location        *string                         `json:"location"`
	Priority        *int32                          `json:"priority"`
	RecurrenceDay   *int32                          `json:"recurrenceDay"`
	RecurrenceMonth *int32                          `json:"recurrenceMonth"`
	RecurrenceMode  database.REMINDERRECURRENCEMODE `json:"recurrenceMode"`
	Status          database.REMINDERSTATUS         `json:"status"`
}

func (u *UpdateReminder) ToUpdateReminderDB() database.UpdateReminderParams {
	return database.UpdateReminderParams{
		ID: u.ID,
		Status: database.NullREMINDERSTATUS{
			Valid:          true,
			REMINDERSTATUS: u.Status,
		},
		Title: types.StrN(u.Title),
		StartAt: pgtype.Timestamptz{
			Time:  u.StartAt,
			Valid: !u.StartAt.IsZero(),
		},
		EndAt: pgtype.Timestamptz{
			Time:  u.EndAt,
			Valid: !u.EndAt.IsZero(),
		},
		Note:            types.StrN(u.Note),
		Location:        types.StrN(u.Location),
		Priority:        types.Int32N(u.Priority),
		RecurrenceDay:   types.Int32N(u.RecurrenceDay),
		RecurrenceMonth: types.Int32N(u.RecurrenceMonth),
		RecurrenceMode: database.NullREMINDERRECURRENCEMODE{
			REMINDERRECURRENCEMODE: u.RecurrenceMode,
			Valid:                  u.RecurrenceMode != "",
		},
	}
}

type GetRemindersQuery struct {
	ResourceTag *string `query:"resourceTag" validate:"omitempty"`
}
