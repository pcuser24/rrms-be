package token

import (
	"time"

	"github.com/google/uuid"
)

type CreateTokenOptions struct {
	TokenID   uuid.UUID
	TokenType TokenType
}

type Maker interface {
	CreateToken(userID uuid.UUID, duration time.Duration, options CreateTokenOptions) (token string, payload *Payload, err error)
	VerifyToken(token string) (payload *Payload, err error)
}
