package dto

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type UpdateProperty struct {
	ID             uuid.UUID
	Name           *string  `json:"name"`
	Area           *float64 `json:"area" validate:"gte:0"`
	NumberOfFloors *int32   `json:"number_of_floors" validate:"gte:0"`
	YearBuilt      *int32   `json:"year_built" validate:"gte:0"`
	Orientation    *string  `json:"orientation" validate:"oneof:n,s,e,w,ne,nw,se,sw"`
	FullAddress    *string  `json:"full_address"`
	District       *string  `json:"district"`
	City           *string  `json:"city"`
	Lat            *float64 `json:"lat"`
	Lng            *float64 `json:"lng"`
}

func (u *UpdateProperty) ToUpdatePropertyDB() *database.UpdatePropertyParams {
	up := database.UpdatePropertyParams{
		ID: u.ID,
	}
	if u.Name != nil {
		up.Name = sql.NullString{String: *u.Name, Valid: true}
	}
	if u.Area != nil {
		up.Area = sql.NullFloat64{Float64: *u.Area, Valid: true}
	}
	if u.NumberOfFloors != nil {
		up.NumberOfFloors = sql.NullInt32{Int32: *u.NumberOfFloors, Valid: true}
	}
	if u.YearBuilt != nil {
		up.YearBuilt = sql.NullInt32{Int32: *u.YearBuilt, Valid: true}
	}
	if u.Orientation != nil {
		up.Orientation = sql.NullString{String: *u.Orientation, Valid: true}
	}
	if u.FullAddress != nil {
		up.FullAddress = sql.NullString{String: *u.FullAddress, Valid: true}
	}
	if u.District != nil {
		up.District = sql.NullString{String: *u.District, Valid: true}
	}
	if u.City != nil {
		up.City = sql.NullString{String: *u.City, Valid: true}
	}
	if u.Lat != nil {
		up.Lat = sql.NullFloat64{Float64: *u.Lat, Valid: true}
	}
	if u.Lng != nil {
		up.Lng = sql.NullFloat64{Float64: *u.Lng, Valid: true}
	}
	return &up
}
