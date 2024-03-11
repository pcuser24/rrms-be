package dto

import "github.com/google/uuid"

type CreateListingPayment struct {
	UserId       uuid.UUID `json:"userId" validate:"required,uuid4"`
	ListingId    uuid.UUID `json:"listingId" validate:"required,uuid4"`
	Priority     int       `json:"priority" validate:"required"`
	PostDuration int       `json:"postDuration" validate:"required"`
}
