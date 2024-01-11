package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type UnitModel struct {
	ID                  uuid.UUID          `json:"id"`
	PropertyID          uuid.UUID          `json:"property_id"`
	Name                string             `json:"name"`
	Area                float32            `json:"area"`
	Floor               *int32             `json:"floor"`
	Price               *int64             `json:"price"`
	NumberOfLivingRooms *int32             `json:"number_of_living_rooms"`
	NumberOfBedrooms    *int32             `json:"number_of_bedrooms"`
	NumberOfBathrooms   *int32             `json:"number_of_bathrooms"`
	NumberOfToilets     *int32             `json:"number_of_toilets"`
	NumberOfKitchens    *int32             `json:"number_of_kitchens"`
	NumberOfBalconies   *int32             `json:"number_of_balconies"`
	Type                database.UNITTYPE  `json:"type"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
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
		Price:               types.PNInt64(u.Price),
		NumberOfLivingRooms: types.PNInt32(u.NumberOfLivingRooms),
		NumberOfBedrooms:    types.PNInt32(u.NumberOfBedrooms),
		NumberOfBathrooms:   types.PNInt32(u.NumberOfBathrooms),
		NumberOfToilets:     types.PNInt32(u.NumberOfToilets),
		NumberOfKitchens:    types.PNInt32(u.NumberOfKitchens),
		NumberOfBalconies:   types.PNInt32(u.NumberOfBalconies),
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}
}
