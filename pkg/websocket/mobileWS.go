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

var topicSLM = make(map[string]map[string][]*websocket.Conn)
var muM sync.Mutex

func StaffidPerWebSocket() fiber.Handler {
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
		muM.Lock()
		if _, ok := topicSLM[topic]; !ok {
			topicSLM[topic] = make(map[string][]*websocket.Conn)
		}
		topicSLM[topic][staffID] = append(topicSLM[topic][staffID], c)

		muM.Unlock()

		// Keep connection alive
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}

		// Remove connection on close
		muM.Lock()
		conns := topicSLM[topic][staffID]
		for i, conn := range conns {
			if conn == c {
				topicSLM[topic][staffID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}

		// Clean up empty slices/maps
		if len(topicSLM[topic][staffID]) == 0 {
			delete(topicSLM[topic], staffID)
		}
		if len(topicSLM[topic]) == 0 {
			delete(topicSLM, topic)
		}
		muM.Unlock()
	})
}

func MPublish(staffID string, data any, topic string) {

	topicCH := topicSLM[topic]
	topicConn := topicCH[staffID]
	fmt.Println("Publishing to topic:", topic, "for staffID:", staffID)
	fmt.Println("Notifying", len(topicConn), "connections for staffID", staffID)

	muM.Lock()
	for _, conn := range topicConn {
		conn.WriteJSON(fiber.Map{
			"topic": topic,
			"type":  staffID, // e.g. "201008-03206"
			"data":  data,
		})
	}
	muM.Unlock()
}
