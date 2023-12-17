package dto

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type CreateUnitAmenity struct {
	AmenityID   int64   `json:"amenityId" validate:"required"`
	Description *string `json:"description"`
}

type CreateUnitMedia struct {
	Url         string             `json:"url" validate:"required,url"`
	Type        database.MEDIATYPE `json:"type" validate:"required,oneof=IMAGE VIDEO"`
	Description *string            `json:"description"`
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
	p := &database.CreateUnitParams{
		PropertyID: cu.PropertyID,
		Area:       cu.Area,
		Type:       cu.Type,
	}
	if cu.Name != nil {
		p.Name = sql.NullString{
			String: *cu.Name,
			Valid:  true,
		}
	}
	if cu.Floor != nil {
		p.Floor = sql.NullInt32{
			Int32: *cu.Floor,
			Valid: true,
		}
	}
	if cu.NumberOfLivingRooms != nil {
		p.NumberOfLivingRooms = sql.NullInt32{
			Int32: *cu.NumberOfLivingRooms,
			Valid: true,
		}
	}
	if cu.NumberOfBedrooms != nil {
		p.NumberOfBedrooms = sql.NullInt32{
			Int32: *cu.NumberOfBedrooms,
			Valid: true,
		}
	}
	if cu.NumberOfBathrooms != nil {
		p.NumberOfBathrooms = sql.NullInt32{
			Int32: *cu.NumberOfBathrooms,
			Valid: true,
		}
	}
	if cu.NumberOfToilets != nil {
		p.NumberOfToilets = sql.NullInt32{
			Int32: *cu.NumberOfToilets,
			Valid: true,
		}
	}
	if cu.NumberOfKitchens != nil {
		p.NumberOfKitchens = sql.NullInt32{
			Int32: *cu.NumberOfKitchens,
			Valid: true,
		}
	}
	if cu.NumberOfBalconies != nil {
		p.NumberOfBalconies = sql.NullInt32{
			Int32: *cu.NumberOfBalconies,
			Valid: true,
		}
	}
	return p
}
