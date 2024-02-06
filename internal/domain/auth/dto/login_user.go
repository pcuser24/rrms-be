package dto

import (
	"time"

	"github.com/google/uuid"
)

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LoginUserRes struct {
	SessionID    uuid.UUID    `json:"sessionId"    `
	AccessToken  string       `json:"accessToken"  `
	AccessExp    time.Time    `json:"accessExp"    `
	RefreshToken string       `json:"refreshToken" `
	RefreshExp   time.Time    `json:"refreshExp"   `
	User         UserResponse `json:"user"         `
}
