package model

import (
	"time"

	"github.com/google/uuid"
)

type PropertyManagerModel struct {
	PropertyID uuid.UUID `json:"propertyId"`
	ManagerID  uuid.UUID `json:"managerId"`
	Role       string    `json:"role"`
}

type NewPropertyManagerRequest struct {
	ID         int64     `json:"id"`
	CreatorID  uuid.UUID `json:"creatorId"`
	PropertyID uuid.UUID `json:"propertyId"`
	UserID     uuid.UUID `json:"userId"`
	Email      string    `json:"email"`
	Approved   bool      `json:"approved"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
