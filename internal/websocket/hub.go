package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"messenger-project/internal/models"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

type Hub struct {
	clients    map[string]*Client
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.username] = client
			h.mutex.Unlock()
			log.Printf("User %s connected to websocket", client.username)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.username]; ok {
				delete(h.clients, client.username)
				close(client.send)
			}
			h.mutex.Unlock()
			log.Printf("User %s disconnected from websocket", client.username)

		case message := <-h.broadcast:
			h.BroadcastToUser(message.ReceiverUsername, message)
		}
	}
}

func (h *Hub) BroadcastToUser(username string, message models.Message) error {
	h.mutex.RLock()
	client, exists := h.clients[username]
	h.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("user %s is not connected", username)
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message for websocket: %v", err)
		return fmt.Errorf("marshal message: %w", err)
	}

	select {
	case client.send <- data:
		return nil
	default:
		close(client.send)
		h.mutex.Lock()
		delete(h.clients, username)
		h.mutex.Unlock()
		return fmt.Errorf("client %s send buffer full or disconnected", username)
	}
}

func (h *Hub) ServeWebSocket(w http.ResponseWriter, r *http.Request, username string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
	}

	h.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}
