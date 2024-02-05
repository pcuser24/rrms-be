package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreateSession struct {
	ID           uuid.UUID `json:"id"`
	UserId       uuid.UUID `json:"userId"`
	SessionToken string    `json:"sessionToken"`
	Expires      time.Time `json:"expires"`
	UserAgent    []byte    `json:"userAgent"`
	ClientIp     string    `json:"clientIp"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (d *CreateSession) ToCreateSessionParams() *database.CreateSessionParams {
	return &database.CreateSessionParams{
		ID:           d.ID,
		Userid:       d.UserId,
		Sessiontoken: d.SessionToken,
		Expires:      d.Expires,
		ClientIp:     types.StrN(&d.ClientIp),
		CreatedAt:    d.CreatedAt,
		UserAgent:    types.StrN(types.Ptr[string](string(d.UserAgent))),
	}
}
