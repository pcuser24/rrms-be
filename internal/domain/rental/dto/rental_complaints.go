package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type PreCreateRentalComplaintMedia struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Size int64  `json:"size" validate:"required,gt=0"`
	Type string `json:"type" validate:"required"`
	Url  string `json:"url"`
}

type PreCreateRentalComplaint struct {
	Media []PreCreateRentalComplaintMedia `json:"media" validate:"dive"`
}

type CreateRentalComplaint struct {
	RentalID   int64                        `json:"rentalId" validate:"required"`
	CreatorID  uuid.UUID                    `json:"creatorId" validate:"required"`
	Title      string                       `json:"title" validate:"required"`
	Content    string                       `json:"content" validate:"required"`
	Suggestion *string                      `json:"suggestion" validate:"omitempty"`
	Media      []string                     `json:"media" validate:"omitempty"`
	OccurredAt time.Time                    `json:"occurredAt" validate:"required"`
	Type       database.RENTALCOMPLAINTTYPE `json:"type" validate:"required"`
}

func (c *CreateRentalComplaint) ToCreateRentalComplaintDB() database.CreateRentalComplaintParams {
	return database.CreateRentalComplaintParams{
		RentalID:   c.RentalID,
		CreatorID:  c.CreatorID,
		Title:      c.Title,
		Content:    c.Content,
		Suggestion: types.StrN(c.Suggestion),
		Media:      c.Media,
		OccurredAt: c.OccurredAt,
		Type:       c.Type,
	}
}

type UpdateRentalComplaint struct {
	Title      *string                        `json:"title"`
	Content    *string                        `json:"content"`
	Suggestion *string                        `json:"suggestion"`
	Media      []string                       `json:"media"`
	ReportAt   time.Time                      `json:"reportAt"`
	Status     database.RENTALCOMPLAINTSTATUS `json:"status"`
	ID         int64                          `json:"id"`
	UserID     uuid.UUID                      `json:"userId"`
}

func (u *UpdateRentalComplaint) ToUpdateRentalComplaintDB() database.UpdateRentalComplaintParams {
	return database.UpdateRentalComplaintParams{
		Title:      types.StrN(u.Title),
		Content:    types.StrN(u.Content),
		Suggestion: types.StrN(u.Suggestion),
		Media:      u.Media,
		OccurredAt: pgtype.Timestamptz{
			Time:  u.ReportAt,
			Valid: !u.ReportAt.IsZero(),
		},
		Status: database.NullRENTALCOMPLAINTSTATUS{
			RENTALCOMPLAINTSTATUS: u.Status,
			Valid:                 u.Status != "",
		},
		ID:     u.ID,
		UserID: u.UserID,
	}
}

type CreateRentalComplaintReply struct {
	ComplaintID int64     `json:"complaintId"`
	ReplierID   uuid.UUID `json:"replierId"`
	Reply       string    `json:"reply"`
	Media       []string  `json:"media"`
}

func (c *CreateRentalComplaintReply) ToCreateRentalComplaintReplyDB() database.CreateRentalComplaintReplyParams {
	return database.CreateRentalComplaintReplyParams{
		ComplaintID: c.ComplaintID,
		ReplierID:   c.ReplierID,
		Reply:       c.Reply,
		Media:       c.Media,
	}
}

type GetRentalComplaintsOfUserQuery struct {
	Limit  int32                          `query:"limit" validate:"omitempty,gte=0"`
	Offset int32                          `query:"offset" validate:"omitempty,gte=0"`
	Status database.RENTALCOMPLAINTSTATUS `query:"status" validate:"omitempty,oneof=PENDING RESOLVED CLOSED"`
}
