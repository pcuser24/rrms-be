package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type ReminderModel struct {
	ID              int64                           `json:"id"`
	CreatorID       uuid.UUID                       `json:"creatorId"`
	Title           string                          `json:"title"`
	StartAt         time.Time                       `json:"startAt"`
	EndAt           time.Time                       `json:"endAt"`
	Note            *string                         `json:"note"`
	Location        string                          `json:"location"`
	Priority        int32                           `json:"priority"`
	Status          database.REMINDERSTATUS         `json:"status"`
	RecurrenceDay   *int32                          `json:"recurrenceDay"`
	RecurrenceMonth *int32                          `json:"recurrenceMonth"`
	RecurrenceMode  database.REMINDERRECURRENCEMODE `json:"recurrenceMode"`
	ResourceTag     string                          `json:"resourceTag"`
	CreatedAt       time.Time                       `json:"createdAt"`
	UpdatedAt       time.Time                       `json:"updatedAt"`

	ReminderMembers []ReminderMemberModel `json:"reminderMembers"`
}

type ReminderMemberModel struct {
	ReminderID int64     `json:"reminderId"`
	UserID     uuid.UUID `json:"userId"`
}

func ToReminderModel(rdb *database.Reminder) *ReminderModel {
	return &ReminderModel{
		ID:              rdb.ID,
		CreatorID:       rdb.CreatorID,
		Title:           rdb.Title,
		StartAt:         rdb.StartAt,
		EndAt:           rdb.EndAt,
		Note:            types.PNStr(rdb.Note),
		Location:        rdb.Location,
		Priority:        rdb.Priority,
		Status:          rdb.Status,
		RecurrenceDay:   types.PNInt32(rdb.RecurrenceDay),
		RecurrenceMonth: types.PNInt32(rdb.RecurrenceMonth),
		RecurrenceMode:  rdb.RecurrenceMode,
		ResourceTag:     rdb.ResourceTag,
		CreatedAt:       rdb.CreatedAt,
		UpdatedAt:       rdb.UpdatedAt,
	}
}
