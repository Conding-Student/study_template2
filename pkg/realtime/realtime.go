package realtime

import (
	"chatbot/pkg/sharedfunctions"
	"encoding/json"
	"log"
	"net/url"
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

// FeatureConnections manages connections for one "feature"
type FeatureConnections struct {
	groups map[string]map[*WSClient]bool // id -> connections
	mu     sync.RWMutex
}

// FeatureList is the global manager holding all features and their connections
type FeatureList struct {
	items map[string]*FeatureConnections // feature name -> FeatureConnections
	mu    sync.RWMutex                   // lock for accessing features map
}

// NewFeatureList creates a new feature list with an empty map
func NewFeatureList() *FeatureList {
	return &FeatureList{
		items: make(map[string]*FeatureConnections),
	}
}

// returns a feature, creating it if needed
func (fl *FeatureList) getOrCreateFeatureconn(feature string) *FeatureConnections {
	fl.mu.Lock()
	defer fl.mu.Unlock() //unlock after checking/creating

	if _, ok := fl.items[feature]; !ok {
		fl.items[feature] = &FeatureConnections{
			groups: make(map[string]map[*WSClient]bool),
		}
		log.Printf("[INIT] Created new connection under specific feature: %s", feature)
	}
	return fl.items[feature]
}

//
// ========== CONNECTION MANAGEMENT ==========
//

// HandleConnection registers a connection under a feature + id
func (h *FeatureList) HandleConnection(c *websocket.Conn, id, feature string) {
	featureConn := h.getOrCreateFeatureconn(feature)

	client := &WSClient{
		Conn: c,
		Send: make(chan []byte, 100), // buffered queue per client
	}

	// Register safely
	featureConn.mu.Lock()
	if _, ok := featureConn.groups[id]; !ok {
		featureConn.groups[id] = make(map[*WSClient]bool)
	}
	featureConn.groups[id][client] = true
	totalIDs := len(featureConn.groups)
	featureConn.mu.Unlock()

	log.Printf("[CONNECT] Feature=%s | ID=%s | Total IDs=%d", feature, id, totalIDs)

	// Start writer goroutine
	go func(cl *WSClient) {
		for msg := range cl.Send {
			if err := cl.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[WRITE ERR] Closing client: %v", err)
				cl.Conn.Close()
				break
			}
		}
	}(client)

	// Cleanup on disconnect
	defer func() {
		featureConn.mu.Lock()
		delete(featureConn.groups[id], client)
		if len(featureConn.groups[id]) == 0 {
			delete(featureConn.groups, id)
		}
		featureConn.mu.Unlock()
		close(client.Send)
		c.Close()
		log.Printf("[DISCONNECT] Feature=%s | ID=%s | Remaining IDs=%d", feature, id, len(featureConn.groups))
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

// Publish sends data to all clients in a feature (optionally filtered by id)
func (h *FeatureList) Publish(id string, feature string, data any) {
	fetchfeatured := h.getOrCreateFeatureconn(feature)

	// Marshal message once
	msg, _ := json.Marshal(data)

	fetchfeatured.mu.RLock()
	defer fetchfeatured.mu.RUnlock()

	if id == "ToAll" {
		for _, clients := range fetchfeatured.groups {
			for client := range clients {
				select {
				case client.Send <- msg:
				default:
					log.Printf("[DROP] Queue full for %v", client.Conn.RemoteAddr())
				}
			}
		}
	} else if clients, ok := fetchfeatured.groups[id]; ok {
		for client := range clients {
			select {
			case client.Send <- msg:
			default:
				log.Printf("[DROP] Queue full for %v", client.Conn.RemoteAddr())
			}
		}
	}
}

//
// ========== AUTHENTICATION ==========
//

// WSAuthMiddleware ensures only valid tokens can connect
func WSAuthMiddleware(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	clientType := strings.ToLower(c.Get("clienttype"))

	var token, id, feature string

	if clientType == "mobile" {
		// Token from headers
		token = c.Get("Authorization")
		if after, ok := strings.CutPrefix(token, "Bearer "); ok {
			token = after
		}
		id = c.Get("id")
		feature = c.Get("feature")
	} else {
		// Web: from query
		token = c.Query("token")
		id = c.Query("id")
		feature = c.Query("feature")
	}

	if token == "" {
		return c.Status(401).SendString("Missing authentication token")
	}

	safeToken := url.QueryEscape(token)
	decodedToken, err := url.QueryUnescape(safeToken)
	if err != nil {
		return c.Status(400).SendString("Invalid token encoding")
	}

	isSuccess, _, _, _, tmessage, err := sharedfunctions.ValidateToken(decodedToken)
	if err != nil || !isSuccess {
		return c.Status(401).SendString(tmessage)
	}

	c.Locals("id", id)
	c.Locals("feature", feature)

	return c.Next()
}

//
// ========== HUB INSTANCE & ROUTES ==========
//

var MainHub = NewFeatureList()

//var Notification = NewFeatureList()

func RealtimeFeatureEndpoint() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		id, _ := c.Locals("id").(string)
		feature, _ := c.Locals("feature").(string)

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
