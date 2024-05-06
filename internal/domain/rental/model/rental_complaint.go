package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type RentalComplaint struct {
	ID         int64                          `json:"id"`
	RentalID   int64                          `json:"rentalId"`
	CreatorID  uuid.UUID                      `json:"creatorId"`
	Title      string                         `json:"title"`
	Content    string                         `json:"content"`
	Suggestion *string                        `json:"suggestion"`
	Media      []string                       `json:"media"`
	OccurredAt time.Time                      `json:"occurredAt"`
	CreatedAt  time.Time                      `json:"createdAt"`
	UpdatedAt  time.Time                      `json:"updatedAt"`
	UpdatedBy  uuid.UUID                      `json:"updatedBy"`
	Type       database.RENTALCOMPLAINTTYPE   `json:"type"`
	Status     database.RENTALCOMPLAINTSTATUS `json:"status"`
}

func ToRentalComplaintModel(rdb *database.RentalComplaint) RentalComplaint {
	return RentalComplaint{
		ID:         rdb.ID,
		RentalID:   rdb.RentalID,
		CreatorID:  rdb.CreatorID,
		Title:      rdb.Title,
		Content:    rdb.Content,
		Suggestion: types.PNStr(rdb.Suggestion),
		Media:      rdb.Media,
		OccurredAt: rdb.OccurredAt,
		CreatedAt:  rdb.CreatedAt,
		UpdatedAt:  rdb.UpdatedAt,
		UpdatedBy:  rdb.UpdatedBy,
		Type:       rdb.Type,
		Status:     rdb.Status,
	}
}

type RentalComplaintReply struct {
	ComplaintID int64     `json:"complaintId"`
	ReplierID   uuid.UUID `json:"replierId"`
	Reply       string    `json:"reply"`
	Media       []string  `json:"media"`
	CreatedAt   time.Time `json:"createdAt"`
}
