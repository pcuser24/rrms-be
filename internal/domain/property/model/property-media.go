package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type PropertyMediaModel struct {
	ID          int64              `json:"id" redis:"id"`
	PropertyID  uuid.UUID          `json:"propertyId" redis:"propertyId"`
	Url         string             `json:"url" redis:"url"`
	Type        database.MEDIATYPE `json:"type" redis:"type"`
	Description *string            `json:"description" redis:"description"`
}

func ToPropertyMediaModel(pm *database.PropertyMedium) PropertyMediaModel {
	return PropertyMediaModel{
		ID:          pm.ID,
		PropertyID:  pm.PropertyID,
		Url:         pm.Url,
		Type:        pm.Type,
		Description: types.PNStr(pm.Description),
	}
}

func (pm PropertyMediaModel) MarshalBinary() (data []byte, err error) {
	return json.Marshal(pm)
}

func (pm *PropertyMediaModel) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, pm)
}
