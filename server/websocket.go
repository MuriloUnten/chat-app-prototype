package main

// TODO Refactor this entire garbage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	websocketBufferSize = 4096
	websocketChannelSize = 256
	maxMsgSize = 4096
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: websocketBufferSize,
	WriteBufferSize: websocketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MsgType string
const (
	ChatMsgType  MsgType = "chat"
	EventMsgType MsgType = "event"
	ErrorMsgType MsgType = "error"
)

type EventType string
const (
	UserJoinedEventType  EventType = "user_joined"
	UserLeftEventType    EventType = "user_left"
	RoomCreatedEventType EventType = "room_created"
	RoomDeletedEventType EventType = "room_deleted"
)

type Message struct {
	Type MsgType         `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ChatMsg struct {
	RoomId  int        `json:"room_id"`
	Sender  UserOutput `json:"sender"`
	Content string     `json:"text"`
}

type EventMsg struct {
	Event  EventType `json:"event"`
	UserId int       `json:"user_id"`
	RoomId int       `json:"room_id"`
}

type ErrorMsg struct {
	Error string `json:"error"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	UserOutput
}

func NewClient(hub *Hub, conn *websocket.Conn, u UserOutput) *Client {
	return &Client{
		hub: hub,
		conn: conn,
		send: make(chan []byte, websocketChannelSize),
		UserOutput: u,
	}
}

type Hub struct {
	clients    map[int]*Client
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client

	// TODO Refactor this ASAP (this is entirely too stupid)
	s *Server
}

func NewHub(s *Server) *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		s: s,
	}
}

func (c *Client) readInboundMsgs() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(pongWait))})

	for {
		msgType, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("websocket error:", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if msgType == websocket.CloseMessage { break }

		var message Message
		err = json.Unmarshal(messageBytes, &message)
		if err != nil {
			c.send <- malformedMessageMsg()
			continue
		}

		if message.Type != ChatMsgType {
			c.send <- malformedMessageMsg()
			continue
		}

		var chat ChatMsg
		err = json.Unmarshal(message.Data, &chat)
		if err != nil {
			c.send <- malformedMessageMsg()
			continue
		}

		c.hub.s.roomsMutex.RLock()
		defer c.hub.s.roomsMutex.RUnlock()
		if c.hub.s.rooms[chat.RoomId] == nil {
			c.send <- roomNotFoundMsg(chat.RoomId)
			continue
		}

		chat.Sender.Id = c.Id
		chat.Sender.Name = c.Name
		message.Data, _ = json.Marshal(chat)

		c.hub.broadcast <- message
	}

	fmt.Println("ws: disconnecting client", c.Id, c.Name)
}

func newEventMsg(eventType EventType, roomId int, userId int) []byte {
	eventMsg := EventMsg{
		Event: eventType,
		RoomId: roomId,
		UserId: userId,
	}
	data, err := json.Marshal(eventMsg)
	if err != nil {
		log.Fatal("error generating constant message:", err)
	}

	msg := Message{
		Type: EventMsgType,
		Data: data,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("error generating constant message:", err)
	}

	return b
}

func newErrorMsg(e string) []byte {
	errMsg := ErrorMsg{
		Error: e,
	}
	data, err := json.Marshal(errMsg)
	if err != nil {
		log.Fatal("error generating constant message:", err)
	}

	msg := Message{
		Type: ErrorMsgType,
		Data: data,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("error generating constant message:", err)
	}

	return b
}

func roomNotFoundMsg(roomId int) []byte {
	return newErrorMsg(fmt.Sprintf("room with id: %d not found", roomId))
}

func malformedMessageMsg() []byte {
	return newErrorMsg("malformed message")
}

func (c *Client) writeOutboundMsgs() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// TODO better format this message
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			err = w.Close()
			if err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}
}

func (h *Hub) Run(s *Server) {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Id] = client

		case client := <-h.unregister:
			if _, ok := h.clients[client.Id]; ok {
				delete(h.clients, client.Id)
				close(client.send)
			}

		case msg := <-h.broadcast:
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Fatal("websocket message json marshalling error:", err)
			}

			switch msg.Type {
			case ChatMsgType:
				var chatMsg ChatMsg
				err := json.Unmarshal(msg.Data, &chatMsg)
				if err != nil {
					continue
				}

				h.BroadcastToRoom(msgBytes, chatMsg.RoomId, chatMsg.Sender.Id)

			case EventMsgType:
				var eventMsg EventMsg
				err := json.Unmarshal(msg.Data, &eventMsg)
				if err != nil {
					continue
				}

				switch eventMsg.Event {
				case UserJoinedEventType, UserLeftEventType:
					h.BroadcastToRoom(msgBytes, eventMsg.RoomId, eventMsg.UserId)

				case RoomCreatedEventType, RoomDeletedEventType:
					h.BroadcastGlobal(msgBytes)

				default:
					continue
				}

			default:
				fmt.Printf("messageType: %s is not handled by broadcast\n", msg.Type)
				continue
			}
		}
	}
}

func (h *Hub) BroadcastToRoom(message []byte, roomId int, senderId int) {
	h.s.roomsMutex.RLock()
	defer h.s.roomsMutex.RUnlock()

	room := h.s.rooms[roomId]
	if room == nil {
		return
	}

	for userId := range room {
		client := h.clients[userId]
		if client == nil {
			continue
		}

		select {
		case client.send <- message:

		default:
			delete(h.clients, userId)
			close(client.send)
		}
	}
}

func (h *Hub) BroadcastGlobal(message []byte) {
	for _, client := range h.clients {
		client.send <- message
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) error {
	var u UserOutput
	userId, err := getIdFromToken(r)
	if err != nil {
		return UserNotAuthenticated()
	}
	u.Id = userId

	q := `SELECT u.name FROM app_user u WHERE u.user_id = $1`
	err = s.db.QueryRow(context.Background(), q, userId).Scan(&u.Name)
	if err != nil {
		return err
	}

	token := r.Header.Get("Sec-Websocket-Protocol")
	localUpgrader := upgrader
	localUpgrader.Subprotocols = []string{token}
	conn, err := localUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("failed to upgrade websocket connection\n", err)
		// must return nil here because upgrader.Upgrade() will have already written the response on error
		// This avoids appending to the response
		return nil
	}

	client := NewClient(s.websocketHub, conn, u)
	client.hub.register <- client

	go client.writeOutboundMsgs()
	go client.readInboundMsgs()

	return nil
}
