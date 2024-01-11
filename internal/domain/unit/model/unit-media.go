package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type UnitMediaModel struct {
	ID          int64              `json:"id"`
	UnitID      uuid.UUID          `json:"unitId"`
	Url         string             `json:"url"`
	Type        database.MEDIATYPE `json:"type"`
	Description *string            `json:"description"`
}

func ToUnitMediaModel(um *database.UnitMedium) *UnitMediaModel {
	return &UnitMediaModel{
		ID:          um.ID,
		UnitID:      um.UnitID,
		Url:         um.Url,
		Type:        um.Type,
		Description: types.PNStr(um.Description),
	}
}
