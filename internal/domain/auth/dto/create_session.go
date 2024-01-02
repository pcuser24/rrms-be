package dto

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type CreateSessionDto struct {
	ID           uuid.UUID `json:"id"`
	UserId       uuid.UUID `json:"userId"`
	SessionToken string    `json:"sessionToken"`
	Expires      time.Time `json:"expires"`
	UserAgent    []byte    `json:"userAgent"`
	ClientIp     string    `json:"clientIp"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (d *CreateSessionDto) ToCreateSessionParams() *database.CreateSessionParams {
	s := &database.CreateSessionParams{
		ID:           d.ID,
		Userid:       d.UserId,
		Sessiontoken: d.SessionToken,
		Expires:      d.Expires,
		ClientIp:     types.StrN(&d.ClientIp),
		CreatedAt:    d.CreatedAt,
	}

	s.UserAgent = sql.NullString{
		String: string(d.UserAgent),
		Valid:  true,
	}

	return s
}
