package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Rooms for match-specific updates
	rooms map[string]map[*Client]bool

	// Mutex for thread-safe operations
	mutex sync.RWMutex
}

// Client represents a websocket client
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// The hub
	hub *Hub

	// Room ID (match ID for cricket scoring)
	roomID string

	// Client ID for identification
	clientID string
}

// Message represents a websocket message
type Message struct {
	Type     string      `json:"type"`
	RoomID   string      `json:"room_id,omitempty"`
	Data     interface{} `json:"data"`
	ClientID string      `json:"client_id,omitempty"`
}

// NewHub creates a new websocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true

	// Add client to room
	if client.roomID != "" {
		if h.rooms[client.roomID] == nil {
			h.rooms[client.roomID] = make(map[*Client]bool)
		}
		h.rooms[client.roomID][client] = true
	}

	log.Printf("Client %s connected to room %s", client.clientID, client.roomID)
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// Remove client from room
		if client.roomID != "" && h.rooms[client.roomID] != nil {
			delete(h.rooms[client.roomID], client)
			if len(h.rooms[client.roomID]) == 0 {
				delete(h.rooms, client.roomID)
			}
		}
	}

	log.Printf("Client %s disconnected from room %s", client.clientID, client.roomID)
}

// broadcastMessage broadcasts a message to all clients
func (h *Hub) broadcastMessage(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// BroadcastToRoom broadcasts a message to all clients in a specific room
func (h *Hub) BroadcastToRoom(roomID string, message interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	room, exists := h.rooms[roomID]
	if !exists {
		return
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range room {
		select {
		case client.send <- messageBytes:
		default:
			close(client.send)
			delete(h.clients, client)
			delete(room, client)
		}
	}
}

// GetRoomClients returns the number of clients in a room
func (h *Hub) GetRoomClients(roomID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	room, exists := h.rooms[roomID]
	if !exists {
		return 0
	}

	return len(room)
}

// GetTotalClients returns the total number of connected clients
func (h *Hub) GetTotalClients() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return len(h.clients)
}

// GetTotalRooms returns the total number of active rooms
func (h *Hub) GetTotalRooms() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return len(h.rooms)
}

// Client methods

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// ServeWS handles websocket requests from clients
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, roomID, clientID string) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 256),
		hub:      h,
		roomID:   roomID,
		clientID: clientID,
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go client.writePump()
	go client.readPump()
}
