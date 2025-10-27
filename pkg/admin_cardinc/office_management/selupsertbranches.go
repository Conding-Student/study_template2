package offices

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

// type SelectBranchesParams struct {
// 	Operation int    `json:"operation"`
// 	Region    string `json:"region"`
// }

type jsonBRequestBody map[string]any

func GetBranches(c *fiber.Ctx) error {
	staffid := c.Params("id") // optional for logging
	getBranchParameters := make(jsonBRequestBody)

	if err := c.BodyParser(&getBranchParameters); err != nil {
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
	result, err := Get_Branch(getBranchParameters)
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

	logs.LOSLogs(c, GetCenterModule, staffid, retCode, getcenter_msg)
	return c.JSON(result)
}

// type UpsertBranchesParams struct {
// 	Operation  int
// 	Region     string
// 	Brcode     string
// 	BranchName string
// 	Active     bool
// 	StaffID    string
// }

func UpsertBranches(c *fiber.Ctx) error {
	staffid := c.Params("id") // Optional, used in logs
	upsertParameters := make(jsonBRequestBody)
	params_select := make(jsonBRequestBody)
	// upsertParameters := new(UpsertBranchesParams)
	// params_select := new(SelectBranchesParams)
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
	result, err := Upsert_Branch(upsertParameters, params_select)
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
