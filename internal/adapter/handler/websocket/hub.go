package websocket

import (
	"log/slog"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/umutaraz/pulseguard/internal/core/domain"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan interface{} // Can be CheckResult or Alert
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		broadcast:  make(chan interface{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
			slog.Info("WS: Client connected", "total", len(h.clients))

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				// conn.Close() // Fiber handles closing usually
				slog.Info("WS: Client disconnected", "total", len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				if err := conn.WriteJSON(message); err != nil {
					slog.Error("WS: Write error, dropping client", "error", err)
					conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}

// BroadcastCheckResult pumps a CheckResult to all connected clients.
func (h *Hub) BroadcastCheckResult(result domain.CheckResult) {
	h.broadcast <- result
}
