package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/pkg/ds/set"
)

type EventHandler func(IncomingEvent, *wsConn) error

type wsConn struct {
	sync.Mutex
	*websocket.Conn
}

func newWSConn(c *websocket.Conn) *wsConn {
	return &wsConn{Conn: c}
}

// protect from concurrent write to websocket connection
func (c *wsConn) WriteMessage(t int, d []byte) error {
	c.Lock()
	defer c.Unlock()
	return c.Conn.WriteMessage(t, d)
}

type WSChatAdapter struct {
	service Service
	sync.RWMutex
	groups         map[GroupIdType]set.Set[*wsConn]
	mapConnToGroup map[*wsConn]GroupIdType
	handlers       map[EventType]EventHandler
	egress         chan OutgoingEvent
}

func (ws *WSChatAdapter) String() string {
	return fmt.Sprintf("groups: %v, mapConnToGroups: %v", ws.groups, ws.mapConnToGroup)
}

func NewWSChatAdapter(s Service) *WSChatAdapter {
	return &WSChatAdapter{
		service:        s,
		groups:         make(map[GroupIdType]set.Set[*wsConn]),
		mapConnToGroup: make(map[*wsConn]GroupIdType),
		handlers:       make(map[EventType]EventHandler),
		egress:         make(chan OutgoingEvent),
	}
}

func (ws *WSChatAdapter) RegisterServer(fibApp *fiber.App, tokenMaker token.Maker) {
	ws.handlers[CHATCREATEMESSAGE] = ws.onCreateMessage
	ws.handlers[CHATDELETEMESSAGE] = ws.onDeleteMessage
	ws.handlers[CHATTYPING] = ws.onTyping

	fibApp.Get("/ws/chat/:id",
		AuthorizedMiddleware(tokenMaker),
		CheckGroupMembership(ws.service),
		websocket.New(func(c *websocket.Conn) {
			gid := c.Locals(GroupIDLocalKey).(GroupIdType)
			conn := newWSConn(c)
			ws.addConn(gid, conn)
			log.Printf("add Conn to group %v, %v\n", gid, ws)

			go ws.receiveMessage(conn, gid)
			ws.sendMessage(conn)
		}))
}

func (ws *WSChatAdapter) addConn(groupId GroupIdType, conn *wsConn) {
	ws.Lock()
	defer ws.Unlock()
	gr, ok := ws.groups[groupId]
	if !ok {
		gr = set.NewSet[*wsConn]()
	}
	gr.Add(conn)
	ws.groups[groupId] = gr
	// ws.mapConnToGroups[conn] = append(ws.mapConnToGroups[conn], groupId)
	ws.mapConnToGroup[conn] = groupId
}

func (ws *WSChatAdapter) addConnToGroup(groupId GroupIdType, conn *wsConn) {
	ws.Lock()
	defer ws.Unlock()
	ws.mapConnToGroup[conn] = groupId
	gr, ok := ws.groups[groupId]
	if !ok {
		gr = set.NewSet[*wsConn]()
	}
	gr.Add(conn)
	ws.groups[groupId] = gr
}

func (ws *WSChatAdapter) getConns(groupId GroupIdType) (set.Set[*wsConn], bool) {
	ws.RLock()
	defer ws.RUnlock()
	gr, ok := ws.groups[groupId]
	return gr, ok
}

func (ws *WSChatAdapter) removeConn(conn *wsConn) {
	ws.Lock()
	defer ws.Unlock()
	grId, ok := ws.mapConnToGroup[conn]
	if !ok {
		return
	}
	gr, ok := ws.groups[grId]
	if !ok {
		return
	}
	gr.Remove(conn)
	if gr.IsEmpty() {
		delete(ws.groups, grId)
	}
	delete(ws.mapConnToGroup, conn)
	conn.Close()
}

func (ws *WSChatAdapter) removeGroup(groupId GroupIdType) {
	ws.Lock()
	defer ws.Unlock()
	gr, ok := ws.groups[groupId]
	if !ok {
		return
	}
	for conn := range gr {
		_, ok := ws.mapConnToGroup[conn]
		if !ok {
			continue
		}
		delete(ws.mapConnToGroup, conn)
	}
	delete(ws.groups, groupId)
}

func (ws *WSChatAdapter) receiveMessage(c *wsConn, groupId GroupIdType) {
	defer func() {
		ws.removeConn(c)
		log.Printf("remove conn from group %v\n ws: %v\n", groupId, ws)
	}()

	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("unexpected closure: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
		var e IncomingEvent
		if err = json.Unmarshal(data, &e); err != nil {
			log.Println("json.Unmarshal: ", err)
			continue
		}
		e.GroupId = groupId
		if err := ws.routeEvent(e, c); err != nil {
			log.Println("route event failed: ", err)
		}
	}
}

func (ws *WSChatAdapter) routeEvent(e IncomingEvent, c *wsConn) error {
	if handler, ok := ws.handlers[e.Type]; ok {
		return handler(e, c)
	}
	return ErrEventNotSupported
}

func (ws *WSChatAdapter) sendMessage(c *wsConn) {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		ws.removeConn(c)
		log.Printf("removed conn\n ws: %v\n", ws)
	}()

	for {
		select {
		case oe, ok := <-ws.egress:
			if !ok {
				log.Println("ws.egress channel closed")
				return
			}
			data, err := json.Marshal(map[string]any{
				"type":    oe.Type,
				"status":  oe.StatusCode,
				"payload": oe.Payload,
			})
			if err != nil {
				log.Println("send message json failed: ", err)
				continue
			}
			err = oe.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("write message failed: ", err)
				ws.removeConn(oe.Conn)
			}
		case <-ticker.C:
			if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("write ping failed: ", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}
