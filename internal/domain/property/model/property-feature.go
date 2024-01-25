package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
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

func ToPropertyFeatureModel(pa *database.PropertyFeature) PropertyFeatureModel {
	return PropertyFeatureModel{
		PropertyID:  pa.PropertyID,
		FeatureID:   pa.FeatureID,
		Description: types.PNStr(pa.Description),
	}
}
