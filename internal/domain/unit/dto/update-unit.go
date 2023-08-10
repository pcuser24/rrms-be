package dto

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type UpdateUnit struct {
	ID                  uuid.UUID `json:"id"`
	Name                *string   `json:"name"`
	Area                *float64  `json:"area"`
	Floor               *int32    `json:"floor"`
	HasBalcony          *bool     `json:"has_balcony"`
	NumberOfLivingRooms *int32    `json:"number_of_living_rooms"`
	NumberOfBedrooms    *int32    `json:"number_of_bedrooms"`
	NumberOfBathrooms   *int32    `json:"number_of_bathrooms"`
	NumberOfToilets     *int32    `json:"number_of_toilets"`
}

func (u *UpdateUnit) ToUpdateUnitDB() *database.UpdateUnitParams {
	p := &database.UpdateUnitParams{
		ID: u.ID,
	}
	if u.Name != nil {
		p.Name = sql.NullString{
			String: *u.Name,
			Valid:  true,
		}
	}
	if u.Area != nil {
		p.Area = sql.NullFloat64{
			Float64: *u.Area,
			Valid:   true,
		}
	}
	if u.Floor != nil {
		p.Floor = sql.NullInt32{
			Int32: *u.Floor,
			Valid: true,
		}
	}
	if u.HasBalcony != nil {
		p.HasBalcony = sql.NullBool{
			Bool:  *u.HasBalcony,
			Valid: true,
		}
	}
	if u.NumberOfLivingRooms != nil {
		p.NumberOfLivingRooms = sql.NullInt32{
			Int32: *u.NumberOfLivingRooms,
			Valid: true,
		}
	}
	if u.NumberOfBedrooms != nil {
		p.NumberOfBedrooms = sql.NullInt32{
			Int32: *u.NumberOfBedrooms,
			Valid: true,
		}
	}
	if u.NumberOfBathrooms != nil {
		p.NumberOfBathrooms = sql.NullInt32{
			Int32: *u.NumberOfBathrooms,
			Valid: true,
		}
	}
	if u.NumberOfToilets != nil {
		p.NumberOfToilets = sql.NullInt32{
			Int32: *u.NumberOfToilets,
			Valid: true,
		}
	}

	return p
}
