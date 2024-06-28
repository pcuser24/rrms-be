package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type ReminderModel struct {
	ID        int64     `json:"id"`
	CreatorID uuid.UUID `json:"creatorId"`
	Title     string    `json:"title"`
	StartAt   time.Time `json:"startAt"`
	EndAt     time.Time `json:"endAt"`
	Note      *string   `json:"note"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// ReminderMembers []ReminderMemberModel `json:"reminderMembers"`
}

func ToReminderModel(rdb *database.Reminder) ReminderModel {
	return ReminderModel{
		ID:        rdb.ID,
		CreatorID: rdb.CreatorID,
		Title:     rdb.Title,
		StartAt:   rdb.StartAt,
		EndAt:     rdb.EndAt,
		Note:      types.PNStr(rdb.Note),
		Location:  rdb.Location,
		CreatedAt: rdb.CreatedAt,
		UpdatedAt: rdb.UpdatedAt,
	}
}
