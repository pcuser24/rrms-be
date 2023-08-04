package dto

import (
	"time"
)

// UserModel without password
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy *string   `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy *string   `json:"updated_by"`
	DeletedF  bool      `json:"deleted_f"`
}
