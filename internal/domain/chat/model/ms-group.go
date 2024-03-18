package model

import (
	"time"

	"github.com/google/uuid"
)

type MsgGroupMember struct {
	GroupID int64     `json:"groupId"`
	UserID  uuid.UUID `json:"userId"`
}

type MsgGroupMemberExtended struct {
	UserID    uuid.UUID `json:"userId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

type MsgGroup struct {
	GroupID   int64     `json:"groupId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy uuid.UUID `json:"createdBy"`

	Members []MsgGroupMember `json:"members"`
}

type MsgGroupExtended struct {
	GroupID   int64     `json:"groupId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy uuid.UUID `json:"createdBy"`

	Members []MsgGroupMemberExtended `json:"members"`
}
