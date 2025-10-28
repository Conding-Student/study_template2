package websocket

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"
	"net/url"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var topicSLW = make(map[string][]*websocket.Conn)
var muW sync.Mutex

func WebPerWebSocket() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		staffID := c.Query("id")
		topic := c.Query("topic")

		if topic == "" {
			c.WriteJSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   "missing staff ID or feature",
					IsSuccess: false,
				},
			})
			return
		}

		// Register connection
		muW.Lock()
		topicSLW[topic] = append(topicSLW[topic], c)
		muW.Unlock()

		// Keep connection alive
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}

		// Remove connection on close
		muW.Lock()
		conns := topicSLW[topic]
		for i, conn := range conns {
			if conn == c {
				topicSLW[topic] = append(conns[:i], conns[i+1:]...)
				break
			}
		}

		// Clean up empty slices/maps
		if len(topicSLW[topic]) == 0 {
			delete(topicSLW, staffID)
		}
		if len(topicSLW[topic]) == 0 {
			delete(topicSLW, topic)
		}
		muW.Unlock()
	})
}

func WPublish(data any, topic string) {

	topicConn := topicSLW[topic]
	fmt.Println("Publishing via Web on topic", len(topicConn), "connections for topic", topic)

	muW.Lock()
	for _, conn := range topicConn {
		conn.WriteJSON(fiber.Map{
			"topic": topic,
			"data":  data,
		})
	}
	muW.Unlock()
}

func WSWebAuthMiddleware(c *fiber.Ctx) error {

	token := c.Query("Authorization")
	decodedToken, err := url.QueryUnescape(token)

	if err != nil {
		return c.Status(400).SendString("Invalid token encoding")
	}

	isValid, _, _, _, msg, err := sharedfunctions.ValidateToken(decodedToken)
	if err != nil || !isValid {
		return c.Status(401).SendString(msg)
	}

	return c.Next()
}
