package model

import "github.com/google/uuid"

type PropertyManagerModel struct {
	PropertyID uuid.UUID `json:"propertyID"`
	ManagerID  uuid.UUID `json:"managerID"`
	Role       string    `json:"role"`
}
