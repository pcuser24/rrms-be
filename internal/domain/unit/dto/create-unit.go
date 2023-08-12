package dto

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type CreateUnitAmenity struct {
	AmenityID   int64   `json:"amenity_id" validate:"required"`
	Description *string `json:"description"`
}

type CreateUnitMedia struct {
	Url  string             `json:"url" validate:"required,url"`
	Type database.MEDIATYPE `json:"type" validate:"required,oneof=IMAGE VIDEO"`
}

type CreateUnit struct {
	PropertyID          uuid.UUID           `json:"property_id" validate:"required,uuid4"`
	Name                *string             `json:"name"`
	Area                float32             `json:"area" validate:"required"`
	Floor               *int32              `json:"floor"`
	HasBalcony          *bool               `json:"has_balcony"`
	NumberOfLivingRooms *int32              `json:"number_of_living_rooms"`
	NumberOfBedrooms    *int32              `json:"number_of_bedrooms"`
	NumberOfBathrooms   *int32              `json:"number_of_bathrooms"`
	NumberOfToilets     *int32              `json:"number_of_toilets"`
	NumberOfKitchens    *int32              `json:"number_of_kitchens"`
	Type                database.UNITTYPE   `json:"type" validate:"required"`
	Amenities           []CreateUnitAmenity `json:"amenities"`
	Medium              []CreateUnitMedia   `json:"medium"`
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
	if cu.HasBalcony != nil {
		p.HasBalcony = sql.NullBool{
			Bool:  *cu.HasBalcony,
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
	return p
}
