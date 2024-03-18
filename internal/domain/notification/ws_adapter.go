package notification

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type useridType = uuid.UUID

type WSNotificationAdapter struct {
	sync.RWMutex
	mapUserIdToConn map[useridType]*websocket.Conn
	egress          chan Notification
}

func NewWSNotificationAdapter() *WSNotificationAdapter {
	return &WSNotificationAdapter{
		mapUserIdToConn: make(map[useridType]*websocket.Conn),
		egress:          make(chan Notification),
	}
}

func (ws *WSNotificationAdapter) Register(fibApp *fiber.App) {
	fibApp.Get("/ws/user/:id", websocket.New(func(c *websocket.Conn) {
		userId := c.Params("id")
		uid, err := uuid.Parse(userId)
		if err != nil {
			c.Close()
			return
		}

		ws.addConn(uid, c)

		go ws.receiveMessage(c, uid)
		ws.sendMessage(c, uid)
	}))
}

func (ws *WSNotificationAdapter) PushMessage(msg Notification) {
	ws.egress <- msg
}

func (ws *WSNotificationAdapter) addConn(userId uuid.UUID, conn *websocket.Conn) {
	ws.Lock()
	defer ws.Unlock()
	ws.mapUserIdToConn[userId] = conn
}

func (ws *WSNotificationAdapter) getConn(userId uuid.UUID) (*websocket.Conn, bool) {
	ws.RLock()
	defer ws.RUnlock()
	conn, err := ws.mapUserIdToConn[userId]
	return conn, err
}

func (ws *WSNotificationAdapter) removeConn(userId uuid.UUID) {
	ws.Lock()
	defer ws.Unlock()
	conn, ok := ws.mapUserIdToConn[userId]
	if !ok {
		return
	}
	conn.Close()
	delete(ws.mapUserIdToConn, userId)
}

func (ws *WSNotificationAdapter) receiveMessage(c *websocket.Conn, userId uuid.UUID) {
	defer ws.removeConn(userId)

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
	}
}

func (ws *WSNotificationAdapter) sendMessage(c *websocket.Conn, userId uuid.UUID) {
	defer ws.removeConn(userId)

	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-ws.egress:
			if !ok {
				return
			}
			conn, ok := ws.getConn(msg.UserId)
			if ok {
				err := conn.WriteMessage(websocket.TextMessage, msg.Payload)
				if err != nil {
					ws.removeConn(msg.UserId)
				}
			}
		case <-ticker.C:
			if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("send ping failed:", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}
