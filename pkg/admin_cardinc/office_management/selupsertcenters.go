package offices

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

// type SelectCentersParams struct {
// 	Operation int    `json:"operation"`
// 	Brcode    string `json:"brcode"`
// 	UnitCode  int    `json:"unitCode"`
// }

var GetCenterModule = "Center Module"

func GetCenters(c *fiber.Ctx) error {
	staffID := c.Params("id")

	getCenterParameters := make(jsonBRequestBody)
	if err := c.BodyParser(&getCenterParameters); err != nil {
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

	// Delegate to query function
	result, err := Get_Center(getCenterParameters)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "failed to fetch centers",
				IsSuccess: false,
				Error:     err,
			},
		})
	}
	retCode := sharedfunctions.GetStringFromMap(result, "retCode")
	getcenter_msg := sharedfunctions.GetStringFromMap(result, "message")

	logs.LOSLogs(c, GetCenterModule, staffID, retCode, getcenter_msg)
	return c.JSON(result)
}

// type UpsertCentersParams struct {
// 	Operation  int    `json:"operation"`
// 	Brcode     string `json:"brcode"`
// 	UnitCode   int    `json:"unitcode"`
// 	CenterCode string `json:"centercode"`
// 	CenterName string `json:"centername"`
// 	StaffID    string `json:"staffid"`
// }

func UpsertCenters(c *fiber.Ctx) error {
	staffid := c.Params("id") // Optional, used in logs
	decision := c.Get("operator")
	upsertParameters := make(jsonBRequestBody)
	params_select := make(jsonBRequestBody)
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

	// Delegate to SQL function
	result, err := Upsert_Center(decision, staffid, upsertParameters, params_select)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to upsert center",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract retCode and message for logging
	retCode := sharedfunctions.GetStringFromMap(result, "retCode")
	msg := sharedfunctions.GetStringFromMap(result, "message")

	logs.LOSLogs(c, GetCenterModule, staffid, retCode, msg)

	return c.JSON(result)
}
