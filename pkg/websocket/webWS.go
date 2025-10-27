package websocket

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var topicSLW = make(map[string]map[string][]*websocket.Conn)
var muW sync.Mutex

func WebPerWebSocket() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		staffID := c.Query("id")
		topic := c.Query("topic")

		if staffID == "" || topic == "" {
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
		if _, ok := topicSLW[topic]; !ok {
			topicSLW[topic] = make(map[string][]*websocket.Conn)
		}
		topicSLW[topic][staffID] = append(topicSLW[topic][staffID], c)

		muW.Unlock()

		// Keep connection alive
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}

		// Remove connection on close
		muW.Lock()
		conns := topicSLW[topic][staffID]
		for i, conn := range conns {
			if conn == c {
				topicSLW[topic][staffID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}

		// Clean up empty slices/maps
		if len(topicSLW[topic][staffID]) == 0 {
			delete(topicSLW[topic], staffID)
		}
		if len(topicSLW[topic]) == 0 {
			delete(topicSLW, topic)
		}
		muW.Unlock()
	})
}

func Publish(staffID string, data any, topic string) {

	topicCH := topicSLW[topic]
	topicConn := topicCH[staffID]
	fmt.Println("Publishing to topic:", topic, "for staffID:", staffID)
	fmt.Println("Notifying", len(topicConn), "connections for staffID", staffID)

	muW.Lock()
	for _, conn := range topicConn {
		conn.WriteJSON(fiber.Map{
			"topic": topic,
			"type":  staffID, // e.g. "201008-03206"
			"data":  data,
		})
	}
	muW.Unlock()
}
