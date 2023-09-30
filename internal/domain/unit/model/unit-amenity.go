package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type UAmenity struct {
	ID      int64  `json:"id"`
	Amenity string `json:"amenity"`
}

type UnitAmenityModel struct {
	UnitID      uuid.UUID `json:"unitID"`
	AmenityID   int64     `json:"amenityID"`
	Description *string   `json:"description"`
}

func ToUnitAmenityModel(ua *database.UnitAmenity) *UnitAmenityModel {
	a := UnitAmenityModel{
		UnitID:    ua.UnitID,
		AmenityID: ua.AmenityID,
	}
	if ua.Description.Valid {
		desc := ua.Description.String
		a.Description = &desc
	}
	return &a
}
