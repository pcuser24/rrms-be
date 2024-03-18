package chat

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/chat/dto"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

func (ws *WSChatAdapter) onCreateMessage(e IncomingEvent, c *websocket.Conn) error {
	conns, ok := ws.getConns(e.GroupId)
	if !ok {
		return nil
	}
	var ie dto.IncomingCreateMessageEvent
	if err := json.Unmarshal(e.Payload, &ie); err != nil {
		return err
	}

	gid := c.Locals(GroupIDLocalKey).(int64)
	tkPayload := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
	ie.From = tkPayload.UserID
	msg, err := ws.service.CreateMessage(gid, &ie)
	if err != nil {
		ws.egress <- OutgoingEvent{
			Conn:       c,
			Type:       CREATEMESSAGE,
			StatusCode: fiber.StatusInternalServerError,
		}
		return err
	}

	var oe dto.OutgoingCreateMessageEvent = *msg
	payload, err := json.Marshal(oe)
	if err != nil {
		return err
	}
	for conn := range conns {
		ws.egress <- OutgoingEvent{
			Conn:       conn,
			Type:       CREATEMESSAGE,
			StatusCode: fiber.StatusCreated,
			Payload:    payload,
		}
	}
	return nil
}

func (ws *WSChatAdapter) onDeleteMessage(e IncomingEvent, c *websocket.Conn) error {
	conns, ok := ws.getConns(e.GroupId)
	if !ok {
		return nil
	}
	var ie dto.IncomingDeleteMessageEvent
	if err := json.Unmarshal(e.Payload, &ie); err != nil {
		return err
	}

	gid := c.Locals(GroupIDLocalKey).(int64)
	tkPayload := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
	ie.DeletedBy = tkPayload.UserID
	rAffected, err := ws.service.DeleteMessage(gid, &ie)
	if err != nil {
		ws.egress <- OutgoingEvent{
			Conn:       c,
			Type:       DELETEMESSAGE,
			StatusCode: fiber.StatusInternalServerError,
		}
		return err
	}
	if rAffected == 0 {
		ws.egress <- OutgoingEvent{
			Conn:       c,
			Type:       DELETEMESSAGE,
			StatusCode: fiber.StatusNotFound,
		}
		return nil
	}

	oe := dto.OutgoingDeleteMessageEvent{
		MessageId: ie.MessageId,
		DeletedBy: ie.DeletedBy,
	}
	payload, err := json.Marshal(oe)
	if err != nil {
		return err
	}
	for conn := range conns {
		ws.egress <- OutgoingEvent{
			Conn:       conn,
			Type:       DELETEMESSAGE,
			StatusCode: fiber.StatusOK,
			Payload:    payload,
		}
	}
	return nil
}

func (ws *WSChatAdapter) onTyping(e IncomingEvent, c *websocket.Conn) error {
	conns, ok := ws.getConns(e.GroupId)
	if !ok {
		return nil
	}

	tkPayload := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
	oe := dto.OutgoingTypingEvent{
		From: tkPayload.UserID,
	}
	payload, err := json.Marshal(oe)
	if err != nil {
		return err
	}
	for conn := range conns {
		ws.egress <- OutgoingEvent{
			Conn:       conn,
			Type:       TYPING,
			StatusCode: fiber.StatusOK,
			Payload:    payload,
		}
	}
	return nil
}
