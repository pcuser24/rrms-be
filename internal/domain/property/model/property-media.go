package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type PropertyMediaModel struct {
	ID         int32              `json:"id"`
	PropertyID uuid.UUID          `json:"property_id"`
	Url        string             `json:"url"`
	Type       database.MEDIATYPE `json:"type"`
}
