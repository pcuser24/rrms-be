package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/chat/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type IncomingCreateMessageEvent struct {
	From    uuid.UUID
	Content string               `json:"content"`
	Type    database.MESSAGETYPE `json:"type"`
}

type OutgoingCreateMessageEvent = model.Message
