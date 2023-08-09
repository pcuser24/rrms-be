package dto

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type CreatePropertyMedia struct {
	// PropertyID uuid.UUID          `json:"property_id" validate:"required,uuid"`
	Url  string             `json:"url" validate:"required,url"`
	Type database.MEDIATYPE `json:"type" validate:"required,oneof=IMAGE VIDEO"`
}

type CreatePropertyAmenity struct {
	// PropertyID  uuid.UUID `json:"property_id" validate:"required,uuid"`
	AmenityID   int64   `json:"amenity" validate:"required"`
	Description *string `json:"description"`
}

type CreatePropertyFeature struct {
	// PropertyID  uuid.UUID `json:"property_id" validate:"required,uuid"`
	FeatureID   int64   `json:"feature" validate:"required"`
	Description *string `json:"description"`
}

type CreatePropertyTag struct {
	// PropertyID uuid.UUID `json:"property_id" validate:"required,uuid"`
	Tag string `json:"tag" validate:"required"`
}

type CreateProperty struct {
	OwnerID        uuid.UUID               `json:"owner_id"`
	Name           *string                 `json:"name"`
	Area           float32                 `json:"area" validate:"required,gt=0"`
	NumberOfFloors *int32                  `json:"number_of_floors"`
	YearBuilt      *int32                  `json:"year_built"`
	Orientation    *string                 `json:"orientation"`
	FullAddress    string                  `json:"full_address" validate:"required"`
	District       string                  `json:"district" validate:"required"`
	City           string                  `json:"city" validate:"required"`
	Lat            float64                 `json:"lat" validate:"required"`
	Lng            float64                 `json:"lng" validate:"required"`
	Type           database.PROPERTYTYPE   `json:"type" validate:"required"`
	Medium         []CreatePropertyMedia   `json:"medium"`
	Amenities      []CreatePropertyAmenity `json:"amenities"`
	Features       []CreatePropertyFeature `json:"features"`
	Tags           []CreatePropertyTag     `json:"tags"`
}

func (c *CreateProperty) ToCreatePropertyDB() *database.CreatePropertyParams {
	p := &database.CreatePropertyParams{
		OwnerID:     c.OwnerID,
		Area:        c.Area,
		FullAddress: c.FullAddress,
		District:    c.District,
		City:        c.City,
		Lat:         c.Lat,
		Lng:         c.Lng,
		Type:        c.Type,
	}
	if c.Name != nil {
		p.Name = sql.NullString{
			Valid:  true,
			String: *c.Name,
		}
	}
	if c.NumberOfFloors != nil {
		p.NumberOfFloors = sql.NullInt32{
			Valid: true,
			Int32: *c.NumberOfFloors,
		}
	}
	if c.YearBuilt != nil {
		p.YearBuilt = sql.NullInt32{
			Valid: true,
			Int32: *c.YearBuilt,
		}
	}
	if c.Orientation != nil {
		p.Orientation = sql.NullString{
			Valid:  true,
			String: *c.Orientation,
		}
	}
	return p
}
