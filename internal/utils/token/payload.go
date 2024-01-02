package token

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	TokenType TokenType `json:"type"`
	UserID    uuid.UUID `json:"sub"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
}

func NewPayload(userId uuid.UUID, duration time.Duration, options CreateTokenOptions) (*Payload, error) {
	var tokenID uuid.UUID = options.TokenID
	if options.TokenType == RefreshToken {
		var err error
		tokenID, err = uuid.NewRandom()
		if err != nil {
			return nil, err
		}
	}

	payload := &Payload{
		ID:        tokenID,
		TokenType: options.TokenType,
		UserID:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ExpiredTokenErr
	}
	return nil
}
