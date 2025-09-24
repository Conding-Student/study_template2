package offices

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

type UpsertClusterParams struct {
	Operation int    `json:"operation"`
	Cluster   string `json:"cluster"`
	StaffID   string `json:"staffID"`
}

var GetClustersModule = "Clusters Module"

func GetClusters(c *fiber.Ctx) error {
	staffid := c.Params("id") // optional for logging

	// Delegate to query function
	result, err := Get_Clusters()
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	retCode := sharedfunctions.GetStringFromMap(result, "retCode")
	message := sharedfunctions.GetStringFromMap(result, "message")
	// Log operation
	logs.LOSLogs(c, GetClustersModule, staffid, retCode, message)

	return c.JSON(result)
}
func UpsertCluster(c *fiber.Ctx) error {
	staffid := c.Params("id") // optional for logging
	upsertParameters := new(UpsertClusterParams)
	if err := c.BodyParser(&upsertParameters); err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	result, err := Upsert_Cluster(staffid, upsertParameters)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}
	retCode := sharedfunctions.GetStringFromMap(result, "retCode")
	message := sharedfunctions.GetStringFromMap(result, "message")

	// Log operation
	logs.LOSLogs(c, GetClustersModule, staffid, retCode, message)
	return c.JSON(result)
}
