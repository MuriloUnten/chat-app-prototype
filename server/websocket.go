package main

// TODO Refactor this entire garbage

import (
	"context"
	"encoding/json"
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
}

type Message struct {
	Type    string         `json:"type"`
	Content MessageContent `json:"content"`
}

type MessageContent struct {
	RoomId int        `json:"room_id"`
	Sender UserOutput `json:"sender"`
	Text   string     `json:"text"`
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if msgType == websocket.CloseMessage { break }

		var message Message
		err = json.Unmarshal(messageBytes, &message)
		if err != nil {
			c.send <- malformedMessage()
			continue
		}

		message.Content.Sender.Id = c.Id
		message.Content.Sender.Name = c.Name

		c.hub.broadcast <- message
	}
}

func malformedMessage() []byte {
	// TODO implement message response for malformed message
	return []byte("malformed message")
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
				// TODO better format thism message
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

type WSRoom struct {
	members map[int]bool
}

func NewWSRoom() *WSRoom {
	return &WSRoom{
		members: make(map[int]bool),
	}
}

type Hub struct {
	clients map[int]*Client
	rooms   map[int]*WSRoom

	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		rooms:      make(map[int]*WSRoom),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run(s *Server) {
	q := `SELECT room_id, user_id FROM room_user`
	rows, err := s.db.Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}

	var roomId, userId int
	for rows.Next() {
		rows.Scan(&roomId, &userId)
		room, ok := h.rooms[roomId]
		if !ok {
			room = NewWSRoom()
		}
		room.members[userId] = true
		h.rooms[roomId] = room
	}

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
			roomId := msg.Content.RoomId
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Fatal("websocket message marshalling error:", err)
			}

			room := h.rooms[roomId]
			if room == nil {
				// TODO would be nice to send a message back telling the client that the room doesn't exist
				break
			}

			for userId := range room.members {
				client := h.clients[userId]
				if client == nil {
					continue
				}

				select {
				case client.send <- msgBytes:

				default:
					delete(h.clients, userId)
					close(client.send)
				}
			}
		}
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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// must return nil here because upgrader.Upgrade() will have already written the response on error
		// This avoids appending to the response
		return nil
	}

	// TODO get all rooms the client previously connected to and update the client info on them

	client := NewClient(s.websocketHub, conn, u)
	client.hub.register <- client

	go client.writeOutboundMsgs()
	go client.readInboundMsgs()

	return nil
}
