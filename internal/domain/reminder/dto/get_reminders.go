package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetRemindersQuery struct {
	CreatorID  uuid.UUID `query:"creatorId" validate:"omitempty"`
	MinStartAt time.Time `query:"minStartAt" validate:"omitempty"`
	MaxStartAt time.Time `query:"maxStartAt" validate:"omitempty"`
}
