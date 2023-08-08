package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

// Elevator, Security camera, Pool, Yard, ...
type PropertyAmenityModel struct {
	PropertyID  uuid.UUID `json:"property_id"`
	Amenity     string    `json:"amenity"`
	Description *string   `json:"description"`
}

func ToPropertyAmenityModel(pa *database.PropertyAmenity) *PropertyAmenityModel {
	a := PropertyAmenityModel{
		PropertyID: pa.PropertyID,
		Amenity:    pa.Amenity,
	}
	if pa.Description.Valid {
		a.Description = &pa.Description.String
	}
	return &a
}

// Security guard, Parking, Gym, ...
type PropertyFeatureModel struct {
	PropertyID  uuid.UUID `json:"property_id"`
	Feature     string    `json:"feature"`
	Description *string   `json:"description"`
}

func ToPropertyFeatureModel(pa *database.PropertyFeature) *PropertyFeatureModel {
	f := PropertyFeatureModel{
		PropertyID: pa.PropertyID,
		Feature:    pa.Feature,
	}
	if pa.Description.Valid {
		f.Description = &pa.Description.String
	}
	return &f
}
