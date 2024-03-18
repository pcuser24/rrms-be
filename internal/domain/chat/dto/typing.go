package dto

import "github.com/google/uuid"

type IncomingTypingEvent struct {
	From string `json:"from"`
}

type OutgoingTypingEvent struct {
	From uuid.UUID `json:"from"`
}
