package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type UnitMediaModel struct {
	ID          int64              `json:"id"`
	UnitID      uuid.UUID          `json:"unitID"`
	Url         string             `json:"url"`
	Type        database.MEDIATYPE `json:"type"`
	Description *string            `json:"description"`
}

func ToUnitMediaModel(um *database.UnitMedia) *UnitMediaModel {
	m := UnitMediaModel{
		ID:     um.ID,
		UnitID: um.UnitID,
		Url:    um.Url,
		Type:   um.Type,
	}
	if um.Description.Valid {
		desc := um.Description.String
		m.Description = &desc
	}
	return &m
}
