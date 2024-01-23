package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UnitModel struct {
	ID                  uuid.UUID          `json:"id"`
	PropertyID          uuid.UUID          `json:"propertyId"`
	Name                string             `json:"name"`
	Area                float32            `json:"area"`
	Floor               *int32             `json:"floor"`
	NumberOfLivingRooms *int32             `json:"numberOfLivingRooms"`
	NumberOfBedrooms    *int32             `json:"numberOfBedrooms"`
	NumberOfBathrooms   *int32             `json:"numberOfBathrooms"`
	NumberOfToilets     *int32             `json:"numberOfToilets"`
	NumberOfKitchens    *int32             `json:"numberOfKitchens"`
	NumberOfBalconies   *int32             `json:"numberOfBalconies"`
	Type                database.UNITTYPE  `json:"type"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
	Amenities           []UnitAmenityModel `json:"amenities"`
	Media               []UnitMediaModel   `json:"media"`
}

func ToUnitModel(u *database.Unit) *UnitModel {
	return &UnitModel{
		ID:                  u.ID,
		PropertyID:          u.PropertyID,
		Name:                u.Name,
		Area:                u.Area,
		Type:                u.Type,
		Floor:               types.PNInt32(u.Floor),
		NumberOfLivingRooms: types.PNInt32(u.NumberOfLivingRooms),
		NumberOfBedrooms:    types.PNInt32(u.NumberOfBedrooms),
		NumberOfBathrooms:   types.PNInt32(u.NumberOfBathrooms),
		NumberOfToilets:     types.PNInt32(u.NumberOfToilets),
		NumberOfKitchens:    types.PNInt32(u.NumberOfKitchens),
		NumberOfBalconies:   types.PNInt32(u.NumberOfBalconies),
		Media:               make([]UnitMediaModel, 0),
		Amenities:           make([]UnitAmenityModel, 0),
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}
}
