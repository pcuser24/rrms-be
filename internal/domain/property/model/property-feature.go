package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

// Security guard, Parking, Gym, ...
type PFeature struct {
	ID      int64  `json:"id"`
	Feature string `json:"feature"`
}

type PropertyFeatureModel struct {
	PropertyID  uuid.UUID `json:"propertyId"`
	FeatureID   int64     `json:"featureId"`
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
