package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const minSecretKeyLen = 32

type JWTMaker struct {
	secreteKey string
}

func NewJWTMaker(secreteKey string) (Maker, error) {
	if len(secreteKey) < minSecretKeyLen {
		return nil, fmt.Errorf("invalid key len: at least %d characters", minSecretKeyLen)
	}
	return &JWTMaker{secreteKey}, nil
}

func (maker *JWTMaker) CreateToken(userID uuid.UUID, duration time.Duration, options CreateTokenOptions) (string, *Payload, error) {
	payload, err := NewPayload(userID, duration, options)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secreteKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secreteKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			payload, _ := jwtToken.Claims.(*Payload)
			return payload, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func (maker *JWTMaker) GetSecreteKey() string {
	return maker.secreteKey
}
