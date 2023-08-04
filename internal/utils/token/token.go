package token

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(userID uuid.UUID, tokenType TokenType, duration time.Duration) (token string, payload *Payload, err error)
	VerifyToken(token string) (payload *Payload, err error)
}
