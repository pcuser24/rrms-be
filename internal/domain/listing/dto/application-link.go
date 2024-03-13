package dto

import "github.com/google/uuid"

type CreateApplicationLink struct {
	ListingId uuid.UUID `json:"listingId" validate:"omitempty"`
	FullName  string    `json:"fullName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone" validate:"required"`
}

type VerifyApplicationLink struct {
	CreateApplicationLink
	Key string `query:"k" validate:"required"`
}
