package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Hub struct {
	clients    map[*WSClient]bool
	rooms      map[string]map[*WSClient]bool
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
}

type WSClient struct {
	hub    *Hub
	conn   *ws.Conn
	send   chan []byte
	rooms  []string
	UserID string
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*WSClient]bool),
		rooms:      make(map[string]map[*WSClient]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			for _, room := range client.rooms {
				if h.rooms[room] == nil {
					h.rooms[room] = make(map[*WSClient]bool)
				}
				h.rooms[room][client] = true
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				for _, room := range client.rooms {
					delete(h.rooms[room], client)
					if len(h.rooms[room]) == 0 {
						delete(h.rooms, room)
					}
				}
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastToRoom(room string, message []byte) {
	if clients, ok := h.rooms[room]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}

func (c *WSClient) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *WSClient) WritePump() {
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
				c.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(ws.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("websocket upgrade error: %v", err)
			return
		}
		userID := c.Query("user_id")
		room := c.Query("room")
		var rooms []string
		if room != "" {
			rooms = append(rooms, room)
		}
		client := &WSClient{
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 256),
			rooms:  rooms,
			UserID: userID,
		}
		hub.register <- client
		go client.WritePump()
		go client.ReadPump()
	}
}
