package dto

import "github.com/google/uuid"

type TaskSendEmailOnNewApplicationPayload struct {
	Username      string    `json:"username"`
	ApplicationId int64     `json:"application_id"`
	ListingId     uuid.UUID `json:"listing_id"`
}
