package chat

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/chat/dto"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

func (ws *WSChatAdapter) onCreateReminder(e IncomingEvent, c *wsConn) error {
	conns, ok := ws.getConns(e.GroupId)
	if !ok {
		return nil
	}
	// var ie dto.IncomingCreateReminderEvent
	// if err := json.Unmarshal(e.Payload, &ie); err != nil {
	// 	return err
	// }

	// var oe dto.OutgoingCreateReminderEvent = *msg
	// payload, err := json.Marshal(oe)
	// if err != nil {
	// 	return err
	// }
	for conn := range conns {
		ws.egress <- OutgoingEvent{
			Conn:       conn,
			Type:       REMINDERCREATE,
			StatusCode: fiber.StatusCreated,
			Payload:    e.Payload,
		}
	}
	return nil
}

func (ws *WSChatAdapter) onUpdateReminder(e IncomingEvent, c *wsConn) error {
	conns, ok := ws.getConns(e.GroupId)
	if !ok {
		return nil
	}
	var ie dto.IncomingUpdateReminderEvent
	if err := json.Unmarshal(e.Payload, &ie); err != nil {
		return err
	}

	tkPayload := c.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
	oe := dto.OutgoingUpdateReminderEvent{
		ID:     ie.ID,
		UserId: tkPayload.UserID,
		Status: ie.Status,
	}
	payload, err := json.Marshal(oe)
	if err != nil {
		return err
	}

	for conn := range conns {
		ws.egress <- OutgoingEvent{
			Conn:       conn,
			Type:       REMINDERUPDATESTATUS,
			StatusCode: fiber.StatusCreated,
			Payload:    payload,
		}
	}
	return nil
}
