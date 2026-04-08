package websocket

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[uuid.UUID]*websocket.Conn
  mu      sync.RWMutex
}

func NewHub() *Hub {
  return &Hub{
    clients: make(map[uuid.UUID]*websocket.Conn),
  }
}

func (h *Hub) AddUser(userID uuid.UUID, conn *websocket.Conn) {
  h.mu.Lock()
  defer h.mu.Unlock()
  h.clients[userID] = conn
}

func (h *Hub) RemoveUser(userID uuid.UUID) {
  h.mu.Lock()
  defer h.mu.Unlock()
  delete(h.clients, userID)
}

func (h *Hub) GetConnection(userID uuid.UUID) (*websocket.Conn, bool) {
  h.mu.RLock()
  defer h.mu.RUnlock()
  conn, ok := h.clients[userID]
  return conn, ok
}

func (h *Hub) SendToUsers(userIDs []uuid.UUID, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, id := range userIDs {
		if conn, ok := h.clients[id]; ok {
			conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (h *Hub) IsOnline(userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}