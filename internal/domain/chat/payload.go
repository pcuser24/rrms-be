package chat

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/contrib/websocket"
)

type EventType string

const (
	CREATEMESSAGE EventType = "CHAT_CREATE_MESSAGE"
	DELETEMESSAGE EventType = "CHAT_DELETE_MESSAGE"
	TYPING        EventType = "CHAT_TYPING"
)

var ErrEventNotSupported = errors.New("event not supported")

type IncomingEvent struct {
	Type    EventType `json:"type"`
	GroupId GroupIdType
	Payload json.RawMessage `json:"payload"`
}

type OutgoingEvent struct {
	Conn       *websocket.Conn
	Type       EventType       `json:"type"`
	StatusCode int             `json:"statusCode"`
	Payload    json.RawMessage `json:"payload"`
}
