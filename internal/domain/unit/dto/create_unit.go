package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreateUnitAmenity struct {
	AmenityID   int64   `json:"amenityId" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
}

type CreateUnitMedia struct {
	Url         string             `json:"url" validate:"required,url"`
	Type        database.MEDIATYPE `json:"type" validate:"required,oneof=IMAGE VIDEO"`
	Description *string            `json:"description" validate:"omitempty"`
}

type CreateUnit struct {
	PropertyID          uuid.UUID           `json:"propertyId" validate:"required,uuid4"`
	Name                *string             `json:"name"`
	Area                float32             `json:"area" validate:"required"`
	Floor               *int32              `json:"floor"`
	NumberOfLivingRooms *int32              `json:"numberOfLivingRooms"`
	NumberOfBedrooms    *int32              `json:"numberOfBedrooms"`
	NumberOfBathrooms   *int32              `json:"numberOfBathrooms"`
	NumberOfToilets     *int32              `json:"numberOfToilets"`
	NumberOfKitchens    *int32              `json:"numberOfKitchens"`
	NumberOfBalconies   *int32              `json:"numberOfBalconies"`
	Type                database.UNITTYPE   `json:"type" validate:"required,oneof=APARTMENT ROOM STUDIO"`
	Amenities           []CreateUnitAmenity `json:"amenities"`
	Media               []CreateUnitMedia   `json:"media"`
}

func (cu *CreateUnit) ToCreateUnitDB() *database.CreateUnitParams {
	return &database.CreateUnitParams{
		PropertyID:          cu.PropertyID,
		Name:                types.StrN(cu.Name),
		Floor:               types.Int32N(cu.Floor),
		NumberOfLivingRooms: types.Int32N(cu.NumberOfLivingRooms),
		NumberOfBedrooms:    types.Int32N(cu.NumberOfBedrooms),
		NumberOfBathrooms:   types.Int32N(cu.NumberOfBathrooms),
		NumberOfToilets:     types.Int32N(cu.NumberOfToilets),
		NumberOfKitchens:    types.Int32N(cu.NumberOfKitchens),
		NumberOfBalconies:   types.Int32N(cu.NumberOfBalconies),
		Area:                cu.Area,
		Type:                cu.Type,
	}
}
