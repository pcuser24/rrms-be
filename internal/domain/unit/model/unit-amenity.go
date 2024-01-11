package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type UAmenity struct {
	ID      int64  `json:"id"`
	Amenity string `json:"amenity"`
}

type UnitAmenityModel struct {
	UnitID      uuid.UUID `json:"unitId"`
	AmenityID   int64     `json:"amenityId"`
	Description *string   `json:"description"`
}

func ToUnitAmenityModel(ua *database.UnitAmenity) *UnitAmenityModel {
	return &UnitAmenityModel{
		UnitID:      ua.UnitID,
		AmenityID:   ua.AmenityID,
		Description: types.PNStr(ua.Description),
	}
}
