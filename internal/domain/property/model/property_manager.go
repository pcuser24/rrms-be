package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PropertyManagerModel struct {
	PropertyID uuid.UUID `json:"propertyId" redis:"propertyId"`
	ManagerID  uuid.UUID `json:"managerId" redis:"managerId"`
	Role       string    `json:"role" redis:"role"`
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

func (pm PropertyManagerModel) MarshalBinary() (data []byte, err error) {
	return json.Marshal(pm)
}

func (pm *PropertyManagerModel) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, pm)
}
