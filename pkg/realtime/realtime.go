package realtime

import (
	"chatbot/pkg/sharedfunctions"
	"net/url"
	"strings"
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
	h.mu.Lock()
	defer h.mu.Unlock()

	if ids, ok := h.groups[feature]; ok {
		for staffID, clients := range ids {
			for conn := range clients {
				if err := conn.WriteJSON(data); err != nil {
					// 🚨 Dead connection: cleanup immediately
					conn.Close()
					delete(clients, conn)
				}
			}

			// If no clients left for this staffID, clean it up
			if len(clients) == 0 {
				delete(ids, staffID)
			}
		}

		// If no IDs left under this feature, clean it up
		if len(ids) == 0 {
			delete(h.groups, feature)
		}
	}
}

// ----------------- AUTH ------------------

func WSAuthMiddleware(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	clientType := c.Get("clienttype")

	var token, id, feature string

	if clientType == "mobile" {
		// Token from headers
		token = c.Get("Authorization")
		if after, ok := strings.CutPrefix(token, "Bearer "); ok {
			token = after
		}

		// Id + feature from headers
		id = c.Get("id")
		feature = c.Get("feature")
	} else {
		// Web: everything from query
		token = c.Query("token")
		id = c.Query("id")
		feature = c.Query("feature")
	}

	if token == "" {
		return c.Status(401).SendString("Missing authentication token")
	}

	// Validate token
	safeToken := url.QueryEscape(token)
	decodedToken, err := url.QueryUnescape(safeToken)
	if err != nil {
		return c.Status(400).SendString("Invalid token encoding")
	}

	isSuccess, _, _, _, tmessage, err := sharedfunctions.ValidateToken(decodedToken)
	if err != nil {
		return c.Status(401).SendString(err.Error())
	}
	if !isSuccess {
		return c.Status(401).SendString(tmessage)
	}

	// Store values so the websocket.Conn handler can read them
	c.Locals("clientType", clientType)
	c.Locals("id", id)
	c.Locals("feature", feature)

	return c.Next()
}

// ----------------- HUB INSTANCE ------------------

var MainHub = NewHub()

// ----------------- REGISTER ROUTES ------------------

func Register(app *fiber.App) {
	// Attach both WSAuthMiddleware and handler.AuthMiddleware etc. as needed
	app.Use("/ws", WSAuthMiddleware)

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {

		// Get values passed from middleware
		id, _ := c.Locals("id").(string)
		feature, _ := c.Locals("feature").(string)

		if feature == "" {
			c.WriteMessage(websocket.TextMessage, []byte("Missing feature"))
			c.Close()
			return
		}
		if id == "" {
			id = "default"
		}

		MainHub.HandleConnection(c, id, feature)
	}))

}
