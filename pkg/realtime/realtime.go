package realtime

import (
	//"chatbot/pkg/sharedfunctions"
	"encoding/json"
	"log"

	//"net/url"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

//
// ========== HUB DESIGN ==========
//

type WSClient struct {
	Conn *websocket.Conn
	Send chan []byte
}

type FeatureConnections struct {
	groups map[string]map[*WSClient]bool // id -> set of clients
	mu     sync.RWMutex
}

type FeatureList struct {
	items map[string]*FeatureConnections // feature -> FeatureConnections
	mu    sync.RWMutex
}

func NewFeatureList() *FeatureList {
	return &FeatureList{
		items: make(map[string]*FeatureConnections),
	}
}

func (fl *FeatureList) getOrCreateFeatureconn(feature string) *FeatureConnections {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	conn, ok := fl.items[feature]
	if !ok {
		conn = &FeatureConnections{groups: make(map[string]map[*WSClient]bool)}
		fl.items[feature] = conn
		log.Printf("[INIT] Created new connection for feature: %s", feature)
	}
	return conn
}

//
// ========== CONNECTION MANAGEMENT ==========
//

func (h *FeatureList) HandleConnection(c *websocket.Conn, id, feature string) {
	fc := h.getOrCreateFeatureconn(feature)

	// Whitelist check
	whitelist, err := GetWhitelist(map[string]any{"featurename": feature})
	if err != nil {
		log.Printf("[ERROR] Failed to get whitelist for feature=%s: %v", feature, err)
		c.WriteMessage(websocket.TextMessage, []byte("Error checking whitelist"))
		c.Close()
		return
	}

	// If whitelist is empty or not found, allow connection
	data, ok := whitelist["data"].([]any)
	if !ok || len(data) == 0 {
		log.Printf("[ALLOW] Feature=%s has no whitelist — open access", feature)
	} else if !isIDInWhitelist(id, whitelist) {
		log.Printf("[DENY] ID=%s not in whitelist for feature=%s", id, feature)
		c.WriteMessage(websocket.TextMessage, []byte("Access denied: ID not in whitelist"))
		c.Close()
		return
	}

	client := &WSClient{
		Conn: c,
		Send: make(chan []byte, 100),
	}

	// Register connection
	fc.mu.Lock()
	if _, ok := fc.groups[id]; !ok {
		fc.groups[id] = make(map[*WSClient]bool)
	}
	fc.groups[id][client] = true
	total := len(fc.groups)
	fc.mu.Unlock()

	log.Printf("[CONNECT] Feature=%s | ID=%s | Total IDs=%d", feature, id, total)

	// Writer goroutine
	go func() {
		for msg := range client.Send {
			if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[WRITE ERR] Closing client (%s): %v", feature, err)
				client.Conn.Close()
				break
			}
		}
	}()

	// Cleanup on disconnect
	defer func() {
		fc.mu.Lock()

		conns := fc.groups[id]
		for c2 := range conns {
			if c2 == client {
				delete(conns, c2)
				break
			}
		}
		if len(conns) == 0 {
			delete(fc.groups, id)
		}

		remaining := len(fc.groups)
		fc.mu.Unlock()

		close(client.Send)
		c.Close()
		log.Printf("[DISCONNECT] Feature=%s | ID=%s | Remaining IDs=%d", feature, id, remaining)
	}()

	// Keep alive
	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

//
// ========== MESSAGE BROADCAST ==========
//

func (h *FeatureList) Publish(id, feature string, data any) {
	fc := h.getOrCreateFeatureconn(feature)

	msg, _ := json.Marshal(data)

	fc.mu.RLock()
	defer fc.mu.RUnlock()

	sendToClients := func(clients map[*WSClient]bool) {
		for client := range clients {
			select {
			case client.Send <- msg:
			default:
				log.Printf("[DROP] Queue full for %v", client.Conn.RemoteAddr())
			}
		}
	}

	if id == "ToAll" {
		for _, clients := range fc.groups {
			sendToClients(clients)
		}
	} else if clients, ok := fc.groups[id]; ok {
		sendToClients(clients)
	}
}

//
// ========== AUTH MIDDLEWARE ==========
//

func WSAuthMiddleware(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	clientType := strings.ToLower(c.Get("clienttype"))

	var token, id, feature string
	if clientType == "mobile" {
		token = strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		id = c.Get("id")
		feature = c.Get("feature")
	} else {
		token = c.Query("token")
		id = c.Query("id")
		feature = c.Query("feature")
	}

	if token == "" {
		return c.Status(401).SendString("Missing authentication token")
	}

	//decodedToken, err := url.QueryUnescape(token)
	//decodedToken, err := url.QueryUnescape(url.QueryEscape(token))

	// if err != nil {
	// 	return c.Status(400).SendString("Invalid token encoding")
	// }

	// isValid, _, _, _, msg, err := sharedfunctions.ValidateToken(decodedToken)
	// if err != nil || !isValid {
	// 	return c.Status(401).SendString(msg)
	// }

	c.Locals("id", id)
	c.Locals("feature", feature)

	return c.Next()
}

//
// ========== HUB INSTANCE & ROUTES ==========
//

var MainHub = NewFeatureList()

func RealtimeFeatureEndpoint() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		id := c.Locals("id").(string)
		feature := c.Locals("feature").(string)

		if feature == "" {
			c.WriteMessage(websocket.TextMessage, []byte("Missing feature"))
			c.Close()
			return
		}
		if id == "" {
			c.WriteMessage(websocket.TextMessage, []byte("Missing staff ID"))
			c.Close()
			return
		}

		MainHub.HandleConnection(c, id, feature)
	})
}
