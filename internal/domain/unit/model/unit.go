package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
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
	um := &UnitModel{
		ID:         u.ID,
		PropertyID: u.PropertyID,
		Name:       u.Name,
		Area:       u.Area,
		Type:       u.Type,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
	if u.Floor.Valid {
		f := u.Floor.Int32
		um.Floor = &f
	}
	if u.NumberOfLivingRooms.Valid {
		n := u.NumberOfLivingRooms.Int32
		um.NumberOfLivingRooms = &n
	}
	if u.NumberOfBedrooms.Valid {
		n := u.NumberOfBedrooms.Int32
		um.NumberOfBedrooms = &n
	}
	if u.NumberOfBathrooms.Valid {
		n := u.NumberOfBathrooms.Int32
		um.NumberOfBathrooms = &n
	}
	if u.NumberOfToilets.Valid {
		n := u.NumberOfToilets.Int32
		um.NumberOfToilets = &n
	}
	if u.NumberOfKitchens.Valid {
		n := u.NumberOfKitchens.Int32
		um.NumberOfKitchens = &n
	}
	if u.NumberOfBalconies.Valid {
		n := u.NumberOfBalconies.Int32
		um.NumberOfBalconies = &n
	}
	return um
}
