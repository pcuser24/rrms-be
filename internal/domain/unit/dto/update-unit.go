package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UpdateUnit struct {
	ID                  uuid.UUID `json:"id"`
	Name                *string   `json:"name"`
	Area                *float32  `json:"area"`
	Floor               *int32    `json:"floor"`
	Price               *int64    `json:"price"`
	NumberOfLivingRooms *int32    `json:"number_of_living_rooms"`
	NumberOfBedrooms    *int32    `json:"number_of_bedrooms"`
	NumberOfBathrooms   *int32    `json:"number_of_bathrooms"`
	NumberOfToilets     *int32    `json:"number_of_toilets"`
	NumberOfKitchens    *int32    `json:"number_of_kitchens"`
	NumberOfBalconies   *int32    `json:"number_of_balconies"`
}

func (u *UpdateUnit) ToUpdateUnitDB() *database.UpdateUnitParams {
	p := &database.UpdateUnitParams{
		ID:                  u.ID,
		Name:                types.StrN(u.Name),
		Area:                types.Float32N(u.Area),
		Floor:               types.Int32N(u.Floor),
		Price:               types.Int64N(u.Price),
		NumberOfLivingRooms: types.Int32N(u.NumberOfLivingRooms),
		NumberOfBedrooms:    types.Int32N(u.NumberOfBedrooms),
		NumberOfBathrooms:   types.Int32N(u.NumberOfBathrooms),
		NumberOfToilets:     types.Int32N(u.NumberOfToilets),
		NumberOfKitchens:    types.Int32N(u.NumberOfKitchens),
		NumberOfBalconies:   types.Int32N(u.NumberOfBalconies),
	}

	return p
}
