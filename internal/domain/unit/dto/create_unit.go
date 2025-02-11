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

type PreCreateUnitMedia struct {
	ID   int64  `json:"id" validate:"required,gt=0"`
	Name string `json:"name" validate:"required"`
	Size int64  `json:"size" validate:"required,gt=0"`
	Type string `json:"type" validate:"required"`
	Url  string `json:"url"`
}

type PreCreateUnit struct {
	Media []PreCreateUnitMedia `json:"media" validate:"dive"`
}

type CreateUnit struct {
	PropertyID          uuid.UUID           `json:"propertyId" validate:"required,uuid4"`
	Name                *string             `json:"name" validate:"omitempty"`
	Area                float32             `json:"area" validate:"required,gte=0"`
	Floor               *int32              `json:"floor" validate:"omitempty,gte=0"`
	NumberOfLivingRooms *int32              `json:"numberOfLivingRooms" validate:"omitempty,gte=0"`
	NumberOfBedrooms    *int32              `json:"numberOfBedrooms" validate:"omitempty,gte=0"`
	NumberOfBathrooms   *int32              `json:"numberOfBathrooms" validate:"omitempty,gte=0"`
	NumberOfToilets     *int32              `json:"numberOfToilets" validate:"omitempty,gte=0"`
	NumberOfKitchens    *int32              `json:"numberOfKitchens" validate:"omitempty,gte=0"`
	NumberOfBalconies   *int32              `json:"numberOfBalconies" validate:"omitempty,gte=0"`
	Type                database.UNITTYPE   `json:"type" validate:"required,oneof=APARTMENT ROOM STUDIO"`
	Amenities           []CreateUnitAmenity `json:"amenities" validate:"dive"`
	Media               []CreateUnitMedia   `json:"media" validate:"dive"`
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
