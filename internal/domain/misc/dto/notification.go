package dto

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type CreateNotificationDevice struct {
	Token    string            `json:"token" validate:"required"`
	Platform database.PLATFORM `json:"platform" validate:"required,oneof=WEB IOS ANDROID"`
}

type CreateNotification struct {
	UserId  uuid.UUID              `json:"userId" validate:"omitempty,uuid4"`
	Title   string                 `json:"title" validate:"required"`
	Content string                 `json:"content" validate:"required"`
	Data    map[string]interface{} `json:"data" validate:"omitempty"`
	Email   bool                   `json:"email" validate:"required"`
	Push    bool                   `json:"push" validate:"required"`
	Sms     bool                   `json:"sms" validate:"required"`
}

func (c *CreateNotification) ToCreateNotificationDB() database.CreateNotificationParams {
	dataBytes, _ := json.Marshal(c.Data)
	return database.CreateNotificationParams{
		UserID: pgtype.UUID{
			Bytes: c.UserId,
			Valid: c.UserId != uuid.Nil,
		},
		Title:   c.Title,
		Content: c.Content,
		Data:    dataBytes,
		Email:   c.Email,
		Push:    c.Push,
		Sms:     c.Sms,
	}
}
