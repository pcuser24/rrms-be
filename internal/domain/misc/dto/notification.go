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

type CreateNotificationTarget struct {
	UserId uuid.UUID `json:"userId" validate:"omitempty,uuid4"`
	Tokens []string  `json:"tokens" validate:"omitempty,dive,required"`
	Emails []string  `json:"emails" validate:"omitempty,dive,email"`
}

type CreateNotification struct {
	Title   string                 `json:"title" validate:"required"`
	Content string                 `json:"content" validate:"required"`
	Data    map[string]interface{} `json:"data" validate:"omitempty"`
	Targets []CreateNotificationTarget
}

func (c *CreateNotification) ToCreateNotificationDB(userId uuid.UUID) database.CreateNotificationParams {
	dataBytes, _ := json.Marshal(c.Data)
	return database.CreateNotificationParams{
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: userId != uuid.Nil,
		},
		Title:   c.Title,
		Content: c.Content,
		Data:    dataBytes,
	}
}
