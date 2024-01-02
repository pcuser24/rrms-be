package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type SessionModel struct {
	ID           uuid.UUID `json:"id"`
	SessionToken string    `json:"sessionToken"`
	UserId       uuid.UUID `json:"userId"`
	Expires      time.Time `json:"expires"`
	UserAgent    *string   `json:"user_agent"`
	ClientIp     *string   `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	CreatedAt    time.Time `json:"created_at"`
}

func ToSessionModel(db *database.Session) *SessionModel {
	m := &SessionModel{
		ID:           db.ID,
		SessionToken: db.SessionToken,
		UserId:       db.UserId,
		Expires:      db.Expires,
		IsBlocked:    db.IsBlocked,
		CreatedAt:    db.CreatedAt,
	}
	if db.UserAgent.Valid {
		m.UserAgent = &db.UserAgent.String
	}
	if db.ClientIp.Valid {
		m.ClientIp = &db.ClientIp.String
	}
	return m
}
