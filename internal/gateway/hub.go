package gateway

import (
	"context"
	"sync"

	"github.com/coder/websocket"
)

type Hub struct {
	mu          sync.RWMutex
	connections map[string]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]map[*websocket.Conn]struct{}),
	}
}

func (h *Hub) Add(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[userID] == nil {
		h.connections[userID] = make(map[*websocket.Conn]struct{})
	}

	h.connections[userID][conn] = struct{}{}
}

func (h *Hub) Remove(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	connections := h.connections[userID]
	delete(connections, conn)

	if len(connections) == 0 {
		delete(h.connections, userID)
	}
}

func (h *Hub) Connections(userID string) []*websocket.Conn {
	h.mu.RLock()
	defer h.mu.RUnlock()

	connections := h.connections[userID]
	result := make([]*websocket.Conn, 0, len(connections))

	for conn := range connections {
		result = append(result, conn)
	}

	return result
}

func (h *Hub) SendToUser(
	ctx context.Context,
	userID string,
	messageType websocket.MessageType,
	data []byte,
) {
	for _, conn := range h.Connections(userID) {
		if err := conn.Write(ctx, messageType, data); err != nil {
			_ = conn.Close(
				websocket.StatusInternalError,
				"failed to deliver message",
			)
		}
	}
}
