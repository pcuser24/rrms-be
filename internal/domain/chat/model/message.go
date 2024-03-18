package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Message struct {
	ID        int64                  `json:"id"`
	GroupID   int64                  `json:"groupId"`
	FromUser  uuid.UUID              `json:"fromUser"`
	Content   string                 `json:"content"`
	Status    database.MESSAGESTATUS `json:"status"`
	Type      database.MESSAGETYPE   `json:"type"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}
