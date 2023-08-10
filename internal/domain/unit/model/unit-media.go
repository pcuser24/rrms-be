package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type UnitMediaModel struct {
	ID     int64              `json:"id"`
	UnitID uuid.UUID          `json:"unit_id"`
	Url    string             `json:"url"`
	Type   database.MEDIATYPE `json:"type"`
}
