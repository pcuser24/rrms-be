package model

import "github.com/google/uuid"

type ListingUnitModel struct {
	ListingID uuid.UUID `json:"listingID"`
	UnitID    uuid.UUID `json:"unitID"`
}
