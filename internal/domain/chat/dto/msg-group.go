package dto

import "github.com/google/uuid"

type CreateMsgGroupMember struct {
	UserId uuid.UUID `json:"userId"`
}

type CreateMsgGroup struct {
	Name    string                 `json:"name"`
	Members []CreateMsgGroupMember `json:"members"`
}
