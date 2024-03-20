package chat

import (
	"encoding/json"
	"errors"
)

type EventType string

const (
	CHATCREATEMESSAGE    EventType = "CHAT_CREATE_MESSAGE"
	CHATDELETEMESSAGE    EventType = "CHAT_DELETE_MESSAGE"
	CHATTYPING           EventType = "CHAT_TYPING"
	REMINDERCREATE       EventType = "REMINDER_CREATE"
	REMINDERUPDATESTATUS EventType = "REMINDER_UPDATE_STATUS"
)

var ErrEventNotSupported = errors.New("event not supported")

type IncomingEvent struct {
	Type    EventType `json:"type"`
	GroupId GroupIdType
	Payload json.RawMessage `json:"payload"`
}

type OutgoingEvent struct {
	Conn       *wsConn
	Type       EventType       `json:"type"`
	StatusCode int             `json:"statusCode"`
	Payload    json.RawMessage `json:"payload"`
}
