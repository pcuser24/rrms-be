package dto

import "github.com/google/uuid"

type IncomingDeleteMessageEvent struct {
	DeletedBy uuid.UUID `json:"deletedBy"`
	MessageId int64     `json:"messageId"`
}

type OutgoingDeleteMessageEvent = IncomingDeleteMessageEvent
