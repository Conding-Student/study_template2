package offices

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

type SelectRegionsParams struct {
	SelectOption int    `json:"operation"`
	Cluster      string `json:"cluster"`
}

var GetRegionModule = "Region Module"

func GetRegions(c *fiber.Ctx) error {
	staffID := c.Params("id") // optional for logging
	getRegionParameters := make(map[string]any)

	// ✅ Parse body into struct
	if err := c.BodyParser(&getRegionParameters); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}
	// Delegate to query function
	result, err := Get_Region(getRegionParameters)
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
	logs.LOSLogs(c, GetRegionModule, staffID, retCode, message)

	return c.JSON(result)
}

type UpsertRegionParams struct {
	Operation int    `json:"operation"`
	Cluster   string `json:"cluster"`
	Region    string `json:"region"`
	StaffID   string `json:"staffID"`
}

func UpsertRegion(c *fiber.Ctx) error {
	staffid := c.Params("id")
	upsertParameters := make(jsonBRequestBody)
	selectparameters := make(jsonBRequestBody)
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

	result, err := Upsert_Region(upsertParameters, selectparameters)
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
	logs.LOSLogs(c, GetRegionModule, staffid, retCode, message)
	return c.JSON(result)
}
