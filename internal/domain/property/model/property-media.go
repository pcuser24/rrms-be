package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type PropertyMediaModel struct {
	ID          int64              `json:"id"`
	PropertyID  uuid.UUID          `json:"propertyID"`
	Url         string             `json:"url"`
	Type        database.MEDIATYPE `json:"type"`
	Description *string            `json:"description"`
}

func ToPropertyMediaModel(pm *database.PropertyMedia) *PropertyMediaModel {
	m := PropertyMediaModel{
		ID:         pm.ID,
		PropertyID: pm.PropertyID,
		Url:        pm.Url,
		Type:       pm.Type,
	}
	if pm.Description.Valid {
		desc := pm.Description.String
		m.Description = &desc
	}
	return &m
}
