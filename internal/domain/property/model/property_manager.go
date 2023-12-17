package model

import "github.com/google/uuid"

type PropertyManagerModel struct {
	PropertyID uuid.UUID `json:"propertyId"`
	ManagerID  uuid.UUID `json:"managerId"`
	Role       string    `json:"role"`
}
