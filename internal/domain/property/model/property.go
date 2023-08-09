package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type PropertyTagModel struct {
	ID         int64     `json:"id"`
	PropertyID uuid.UUID `json:"property_id"`
	Tag        string    `json:"tag"`
}

type PropertyModel struct {
	ID             uuid.UUID `json:"id"`
	OwnerID        uuid.UUID `json:"owner_id"`
	Name           *string   `json:"name"`
	Area           float32   `json:"area"`
	NumberOfFloors *int32    `json:"number_of_floors"`
	YearBuilt      *int32    `json:"year_built"`
	// n,s,w,e,nw,ne,sw,se
	Orientation *string                `json:"orientation"`
	FullAddress string                 `json:"full_address"`
	District    string                 `json:"district"`
	City        string                 `json:"city"`
	Lat         float64                `json:"lat"`
	Lng         float64                `json:"lng"`
	Type        database.PROPERTYTYPE  `json:"type"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Amenities   []PropertyAmenityModel `json:"amenities"`
	Features    []PropertyFeatureModel `json:"features"`
	Medium      []PropertyMediaModel   `json:"medium"`
	Tags        []PropertyTagModel     `json:"tags"`
}

func ToPropertyModel(p *database.Property) *PropertyModel {
	m := &PropertyModel{
		ID:          p.ID,
		OwnerID:     p.OwnerID,
		Name:        &p.Name,
		Area:        p.Area,
		FullAddress: p.FullAddress,
		District:    p.District,
		City:        p.City,
		Lat:         p.Lat,
		Lng:         p.Lng,
		Type:        p.Type,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if p.NumberOfFloors.Valid {
		n := p.NumberOfFloors.Int32
		m.NumberOfFloors = &n
	}
	if p.YearBuilt.Valid {
		y := p.YearBuilt.Int32
		m.YearBuilt = &y
	}
	if p.Orientation.Valid {
		o := p.Orientation.String
		m.Orientation = &o
	}

	return m
}
