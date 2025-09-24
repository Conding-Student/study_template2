// package realtime

// import (
// 	"chatbot/pkg/sharedfunctions"
// 	"net/url"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/websocket/v2"
// )

// // Hub represents one WebSocket "room"/topic (e.g., articles, trivia, etc.)
// type Hub struct {
// 	clients   map[*websocket.Conn]bool
// 	broadcast chan map[string]any
// }

// // NewHub creates and starts a new Hub
// func NewHub() *Hub {
// 	h := &Hub{
// 		clients:   make(map[*websocket.Conn]bool),
// 		broadcast: make(chan map[string]any),
// 	}
// 	go h.startBroadcaster()
// 	return h
// }

// // startBroadcaster listens for messages and pushes them to all clients
// func (h *Hub) startBroadcaster() {
// 	for {
// 		msg := <-h.broadcast
// 		for client := range h.clients {
// 			err := client.WriteJSON(msg)
// 			if err != nil {
// 				client.Close()
// 				delete(h.clients, client)
// 			}
// 		}
// 	}
// }

// // HandleConnection adds/removes clients and listens for disconnects
// func (h *Hub) HandleConnection(c *websocket.Conn) {
// 	h.clients[c] = true
// 	defer func() {
// 		delete(h.clients, c)
// 		c.Close()
// 	}()

// 	for {
// 		// We don't use the client message for now, just keep connection alive
// 		if _, _, err := c.ReadMessage(); err != nil {
// 			delete(h.clients, c)
// 			break
// 		}
// 	}
// }

// // Publish sends data to all clients in this hub
// func (h *Hub) Publish(data map[string]any) {
// 	h.broadcast <- data
// }

// // WebSocket authentication middleware
// func WSAuthMiddleware(c *fiber.Ctx) error {
// 	// Ensure it's a WebSocket upgrade request
// 	if !websocket.IsWebSocketUpgrade(c) {
// 		return fiber.ErrUpgradeRequired
// 	}

// 	// Extract token from query string (?token=123)
// 	token := c.Query("token")

// 	safeToken := url.QueryEscape(token)
// 	decodedToken, err := url.QueryUnescape(safeToken)
// 	if err != nil {
// 		return c.Status(400).SendString("Invalid token encoding")
// 	}

// 	if safeToken == "" {
// 		return c.Status(401).SendString("Missing authentication token")
// 	}

// 	isSuccess, _, _, _, tmessage, err := sharedfunctions.ValidateToken(decodedToken)

// 	if err != nil {
// 		return c.Status(401).SendString(err.Error())
// 	}
// 	if !isSuccess {
// 		return c.Status(401).SendString(tmessage)
// 	}

// 	return c.Next()
// }

// // -----------------------------------------//
// //       Create multiple hubs (topics)      //
// // -----------------------------------------//

// var (
// 	ArticlesHub = NewHub()
// 	TriviaHub   = NewHub()
// 	//offices management
// 	UpsertCentersHub = NewHub()
// 	UpsertClusterHub = NewHub()
// 	UpsertRegionHub  = NewHub()
// 	UpsertUnitsHub   = NewHub()
// 	//Mlni user management
// 	MlniStaffHub = NewHub()
// 	//logs
// 	CagabayLogsHub = NewHub()
// )

// // func NotifyLoanUpdate(staffID string)
// // { mu.Lock() conns := staffListeners[staffID] mu.Unlock() if len(conns) == 0 { return }
// // // Get latest data data, _, _, _, _, _, err := GetAllLoans(staffID) if err != nil { return } loans := sharedfunctions.GetListAny(data, "data") // Send to all connected clients of that staffID for _, conn := range conns { conn.WriteJSON(fiber.Map{ "type": staffID, // e.g. "201008-03206" "data": loans, }) } }

// // Register all WebSocket endpoints
// func Register(app *fiber.App) {
// 	// Middleware: only allow WS upgrades
// 	app.Use("/ws", func(c *fiber.Ctx) error {
// 		if websocket.IsWebSocketUpgrade(c) {
// 			return c.Next()
// 		}
// 		return fiber.ErrUpgradeRequired
// 	})

// 	// Each endpoint uses its own hub
// 	//For non-SSL: ws://localhost:PORT/ws/articles
// 	//For SSL: wss://localhost:PORT/ws/articles

// }
package realtime

import (
	"chatbot/pkg/sharedfunctions"
	"net/url"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Hub manages websocket groups by feature (with room for id expansion)
type Hub struct {
	groups map[string]map[string]map[*websocket.Conn]bool // feature -> id -> clients
	mu     sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		groups: make(map[string]map[string]map[*websocket.Conn]bool),
	}
}

// HandleConnection registers a client under feature (+id for possible expansion)
func (h *Hub) HandleConnection(c *websocket.Conn, id string, feature string) {
	h.mu.Lock()
	if _, ok := h.groups[feature]; !ok {
		h.groups[feature] = make(map[string]map[*websocket.Conn]bool)
	}
	if _, ok := h.groups[feature][id]; !ok {
		h.groups[feature][id] = make(map[*websocket.Conn]bool)
	}
	h.groups[feature][id][c] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.groups[feature][id], c)
		if len(h.groups[feature][id]) == 0 {
			delete(h.groups[feature], id)
		}
		if len(h.groups[feature]) == 0 {
			delete(h.groups, feature)
		}
		h.mu.Unlock()
		c.Close()
	}()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

// Publish sends data to all clients subscribed to a feature
// Later: can extend this to also filter by id
// Publish sends data to all clients subscribed to a feature (ignores id for now)
func (h *Hub) Publish(id string, feature string, data any) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if ids, ok := h.groups[feature]; ok {
		for _, clients := range ids {
			for conn := range clients {
				if err := conn.WriteJSON(data); err != nil {
					conn.Close()
					delete(clients, conn)
				}
			}
		}
	}
}

// ----------------- AUTH ------------------

func WSAuthMiddleware(c *fiber.Ctx) error {
	// Ensure it's a WebSocket upgrade request
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	// Extract token from query string (?token=123)
	token := c.Query("token")

	safeToken := url.QueryEscape(token)
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

// ----------------- HUB INSTANCE ------------------

var MainHub = NewHub()

// ----------------- REGISTER ROUTES ------------------

func Register(app *fiber.App) {
	// Attach both WSAuthMiddleware and handler.AuthMiddleware etc. as needed
	app.Use("/ws", WSAuthMiddleware)

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		id := c.Query("id")
		feature := c.Query("feature")

		if feature == "" {
			c.WriteMessage(websocket.TextMessage, []byte("Missing feature"))
			c.Close()
			return
		}
		if id == "" {
			id = "default" // placeholder if no id yet
		}

		MainHub.HandleConnection(c, id, feature)
	}))
}
