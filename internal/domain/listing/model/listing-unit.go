package model

import "github.com/google/uuid"

type ListingUnitModel struct {
	ListingID uuid.UUID `json:"listing_id"`
	UnitID    uuid.UUID `json:"unit_id"`
}
