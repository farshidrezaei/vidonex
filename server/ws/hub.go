// Package ws implements thread-safe WebSocket connections and message broadcasting for live telemetry.
package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true // Allow all origins in local/workstation dev mode
	},
}

// Client represents a connected WebSocket client.
type Client struct {
	hub               *Hub
	socketConnection  *websocket.Conn
	outgoingMessageCh chan []byte
}

// MessageEnvelope wraps outgoing WebSocket events.
type MessageEnvelope struct {
	Type      string      `json:"type"` // e.g. "render_progress", "render_complete", "render_failed"
	Payload   any         `json:"payload"`
	Timestamp int64       `json:"timestamp"`
}

// Hub manages active WebSocket clients and broadcasts events.
type Hub struct {
	clients         map[*Client]bool
	broadcastCh     chan []byte
	registerCh      chan *Client
	unregisterCh    chan *Client
	mutex           sync.RWMutex
	logger          *slog.Logger
}

// NewHub initializes a new WebSocket Hub.
func NewHub(logger *slog.Logger) *Hub {
	if logger == nil {
		logger = slog.Default()
	}
	return &Hub{
		clients:         make(map[*Client]bool),
		broadcastCh:     make(chan []byte, 256),
		registerCh:      make(chan *Client),
		unregisterCh:    make(chan *Client),
		logger:          logger,
	}
}

// Run starts the WebSocket event dispatch loop.
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.registerCh:
			hub.mutex.Lock()
			hub.clients[client] = true
			hub.mutex.Unlock()
			hub.logger.Debug("websocket client connected", "active_clients", len(hub.clients))

		case client := <-hub.unregisterCh:
			hub.mutex.Lock()
			if _, exists := hub.clients[client]; exists {
				delete(hub.clients, client)
				close(client.outgoingMessageCh)
			}
			hub.mutex.Unlock()
			hub.logger.Debug("websocket client disconnected", "active_clients", len(hub.clients))

		case message := <-hub.broadcastCh:
			hub.mutex.RLock()
			for client := range hub.clients {
				select {
				case client.outgoingMessageCh <- message:
				default:
					// Slow client; close and unregister
					close(client.outgoingMessageCh)
					delete(hub.clients, client)
				}
			}
			hub.mutex.RUnlock()
		}
	}
}

// Broadcast sends a typed message payload to all connected clients.
func (hub *Hub) Broadcast(eventType string, payload any) {
	envelope := MessageEnvelope{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now().UnixMilli(),
	}
	bytes, err := json.Marshal(envelope)
	if err != nil {
		hub.logger.Error("failed marshaling websocket message", "error", err)
		return
	}
	hub.broadcastCh <- bytes
}

// HandleWebSocket upgrades HTTP requests to WebSocket connections and registers the client.
func (hub *Hub) HandleWebSocket(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	connection, err := websocketUpgrader.Upgrade(responseWriter, httpRequest, nil)
	if err != nil {
		hub.logger.Error("failed upgrading to websocket", "error", err)
		return
	}

	client := &Client{
		hub:               hub,
		socketConnection:  connection,
		outgoingMessageCh: make(chan []byte, 64),
	}

	hub.registerCh <- client

	go client.writePump()
	go client.readPump()
}

func (client *Client) readPump() {
	defer func() {
		client.hub.unregisterCh <- client
		_ = client.socketConnection.Close()
	}()

	client.socketConnection.SetReadLimit(512 * 1024)
	_ = client.socketConnection.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.socketConnection.SetPongHandler(func(string) error {
		_ = client.socketConnection.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.socketConnection.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		_ = client.socketConnection.Close()
	}()

	for {
		select {
		case message, ok := <-client.outgoingMessageCh:
			_ = client.socketConnection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = client.socketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := client.socketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = writer.Write(message)
			if err := writer.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = client.socketConnection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.socketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
