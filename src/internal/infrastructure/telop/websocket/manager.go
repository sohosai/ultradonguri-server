package websocket

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// websocket の通信の集合
type WebSocketHub struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]bool
	// telopChannel chan entities.TelopMessage
	telopChannel chan WebSocketResponse
}

type WebSocketResponse struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func NewWebSocketHub(bufferSize int) *WebSocketHub {
	return &WebSocketHub{
		conns: make(map[*websocket.Conn]bool),
		// telopChannel: make(chan entities.TelopMessage, bufferSize),
		telopChannel: make(chan WebSocketResponse, bufferSize),
	}
}

func (h *WebSocketHub) AddConnection(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.conns[conn] = true
}

func (h *WebSocketHub) RemoveConnection(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns, conn)
	conn.Close()
}

// func (h *WebSocketHub) PushTelop(telop entities.TelopMessage) {
func (h *WebSocketHub) PushTelop(telop WebSocketResponse) {
	h.telopChannel <- telop
}

func (h *WebSocketHub) StartTelopWebsocketBroadcastWorker() {
	for telop := range h.telopChannel {
		h.mu.Lock()
		conns := make([]*websocket.Conn, 0, len(h.conns))
		for conn := range h.conns {
			conns = append(conns, conn)
		}
		h.mu.Unlock()

		payload, err := json.Marshal(telop)
		if err != nil {
			slog.Error("Marshal error: " + err.Error())
			continue
		}

		for _, conn := range conns {
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			err := conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				slog.Error("Error sending message: " + err.Error())
				h.RemoveConnection(conn)
			}
		}
	}
}
