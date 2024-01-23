package model

import "github.com/google/uuid"

type ListingUnitModel struct {
	ListingID uuid.UUID `json:"listingId"`
	UnitID    uuid.UUID `json:"unitId"`
	Price     int64     `json:"price"`
}
