package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type PropertyMediaModel struct {
	ID          int64              `json:"id"`
	PropertyID  uuid.UUID          `json:"propertyId"`
	Url         string             `json:"url"`
	Type        database.MEDIATYPE `json:"type"`
	Description *string            `json:"description"`
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
