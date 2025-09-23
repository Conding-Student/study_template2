package realtime

import (
	"chatbot/pkg/sharedfunctions"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Hub represents one WebSocket "room"/topic (e.g., articles, trivia, etc.)
type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan map[string]any
}

// NewHub creates and starts a new Hub
func NewHub() *Hub {
	h := &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan map[string]any),
	}
	go h.startBroadcaster()
	return h
}

// startBroadcaster listens for messages and pushes them to all clients
func (h *Hub) startBroadcaster() {
	for {
		msg := <-h.broadcast
		for client := range h.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}

// HandleConnection adds/removes clients and listens for disconnects
func (h *Hub) HandleConnection(c *websocket.Conn) {
	h.clients[c] = true
	defer func() {
		delete(h.clients, c)
		c.Close()
	}()

	for {
		// We don't use the client message for now, just keep connection alive
		if _, _, err := c.ReadMessage(); err != nil {
			delete(h.clients, c)
			break
		}
	}
}

// Publish sends data to all clients in this hub
func (h *Hub) Publish(data map[string]any) {
	h.broadcast <- data
}

// WebSocket authentication middleware
func WSAuthMiddleware(c *fiber.Ctx) error {
	// Ensure it's a WebSocket upgrade request
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	// Extract token from query string (?token=123)
	token := c.Query("token")

	safeToken := url.QueryEscape(token)
	// Fix "+" being converted to spaces
	//token = strings.ReplaceAll(token, " ", "+")
	decodedToken, err := url.QueryUnescape(safeToken)
	if err != nil {
		return c.Status(400).SendString("Invalid token encoding")
	}

	if safeToken == "" {
		return c.Status(401).SendString("Missing authentication token")
	}

	isSuccess, _, _, _, tmessage, err := sharedfunctions.ValidateToken(decodedToken)

	if err != nil {
		return c.Status(401).SendString(err.Error())
	}
	if !isSuccess {
		return c.Status(401).SendString(tmessage)
	}

	return c.Next()
}

// -----------------------------------------//
//       Create multiple hubs (topics)      //
// -----------------------------------------//

var (
	ArticlesHub = NewHub()
	TriviaHub   = NewHub()
	//offices management
	UpsertCentersHub = NewHub()
	UpsertClusterHub = NewHub()
	UpsertRegionHub  = NewHub()
	UpsertUnitsHub   = NewHub()
	//Mlni user management
	MlniStaffHub = NewHub()
	//logs
	CagabayLogsHub = NewHub()
)

// func NotifyLoanUpdate(staffID string)
// { mu.Lock() conns := staffListeners[staffID] mu.Unlock() if len(conns) == 0 { return }
// // Get latest data data, _, _, _, _, _, err := GetAllLoans(staffID) if err != nil { return } loans := sharedfunctions.GetListAny(data, "data") // Send to all connected clients of that staffID for _, conn := range conns { conn.WriteJSON(fiber.Map{ "type": staffID, // e.g. "201008-03206" "data": loans, }) } }

// Register all WebSocket endpoints
func Register(app *fiber.App) {
	// Middleware: only allow WS upgrades
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Each endpoint uses its own hub
	//For non-SSL: ws://localhost:PORT/ws/articles
	//For SSL: wss://localhost:PORT/ws/articles

}
