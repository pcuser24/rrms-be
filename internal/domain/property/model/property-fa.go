package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

// Elevator, Security camera, Pool, Yard, ...
type PAmenity struct {
	ID      int64  `json:"id"`
	Amenity string `json:"amenity"`
}

type PropertyAmenityModel struct {
	PropertyID  uuid.UUID `json:"property_id"`
	AmenityID   int64     `json:"amenity"`
	Description *string   `json:"description"`
}

func ToPropertyAmenityModel(pa *database.PropertyAmenity) *PropertyAmenityModel {
	a := PropertyAmenityModel{
		PropertyID: pa.PropertyID,
		AmenityID:  pa.AmenityID,
	}
	if pa.Description.Valid {
		desc := pa.Description.String
		a.Description = &desc
	}
	return &a
}

// Security guard, Parking, Gym, ...
type PFeature struct {
	ID      int64  `json:"id"`
	Feature string `json:"feature"`
}

type PropertyFeatureModel struct {
	PropertyID  uuid.UUID `json:"property_id"`
	FeatureID   int64     `json:"feature"`
	Description *string   `json:"description"`
}

func ToPropertyFeatureModel(pa *database.PropertyFeature) *PropertyFeatureModel {
	f := PropertyFeatureModel{
		PropertyID: pa.PropertyID,
		FeatureID:  pa.FeatureID,
	}
	if pa.Description.Valid {
		desc := pa.Description.String
		f.Description = &desc
	}
	return &f
}
