package dto

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

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
