package model

import (
	"encoding/json"

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
	PropertyID  uuid.UUID `json:"propertyId" redis:"propertyId"`
	FeatureID   int64     `json:"featureId" redis:"featureId"`
	Description *string   `json:"description" redis:"description"`
}

func ToPropertyFeatureModel(pa *database.PropertyFeature) PropertyFeatureModel {
	return PropertyFeatureModel{
		PropertyID:  pa.PropertyID,
		FeatureID:   pa.FeatureID,
		Description: types.PNStr(pa.Description),
	}
}

func (pf PropertyFeatureModel) MarshalBinary() (data []byte, err error) {
	return json.Marshal(pf)
}

func (pf *PropertyFeatureModel) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, pf)
}
