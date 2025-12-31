package websocket

import (
	"log/slog"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func NewWebSocketHandler(hub *Hub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// Register connection
		hub.register <- c
		defer func() {
			hub.unregister <- c
			c.Close()
		}()

		// Keep connection alive / Listen for incoming (if any)
		// For dashboard, we mostly push OUT, but client might send "subscribe" etc.
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Warn("WS: Unexpected close", "error", err)
				}
				break
			}
		}
	})
}

// Middleware to upgrade HTTP to WS
func UpgradeMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return c.SendStatus(fiber.StatusUpgradeRequired)
}
