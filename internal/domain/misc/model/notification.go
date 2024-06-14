package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type NotificationDevice struct {
	UserID       uuid.UUID         `json:"userId"`
	SessionID    uuid.UUID         `json:"sessionId"`
	Token        string            `json:"token"`
	Platform     database.PLATFORM `json:"platform"`
	LastAccessed time.Time         `json:"lastAccessed"`
	CreatedAt    time.Time         `json:"createdAt"`
}

type Notification struct {
	ID        int64                        `json:"id"`
	UserID    uuid.UUID                    `json:"userId"`
	Title     string                       `json:"title"`
	Content   string                       `json:"content"`
	Data      map[string]interface{}       `json:"data"`
	Target    string                       `json:"target"`
	Channel   database.NOTIFICATIONCHANNEL `json:"channel"`
	Seen      bool                         `json:"seen"`
	CreatedAt time.Time                    `json:"createdAt"`
	UpdatedAt time.Time                    `json:"updatedAt"`
}

func ToNotificationModel(n database.Notification) Notification {
	nm := Notification{
		ID:        n.ID,
		UserID:    n.UserID.Bytes,
		Title:     n.Title,
		Content:   n.Content,
		Seen:      n.Seen,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	if err := json.Unmarshal(n.Data, &nm.Data); err != nil {
		nm.Data = map[string]interface{}{}
	}

	return nm
}
