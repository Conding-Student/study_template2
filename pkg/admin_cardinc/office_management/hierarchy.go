// BY Norman Villegas
package offices

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func GetStaffName(c *fiber.Ctx) error {

	// Parse request JSON
	var reqBody map[string]any
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request!",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request",
				Error:     err,
			},
		})
	}

	var result map[string]any
	result, err := Get_fullname(reqBody)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "failed to full name",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}
func GetStaffByDesignation(c *fiber.Ctx) error {
	var reqBody map[string]any
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request!",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request",
				Error:     err,
			},
		})
	}

	result, err := GetStaffByDesignationDB(reqBody)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Problem connecting to database",
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}

func GetCenterByStaffID(c *fiber.Ctx) error {
	var reqBody map[string]any
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request!",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request",
				Error:     err,
			},
		})
	}

	result, err := GetCenterByStaffIDDB(reqBody)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Problem connecting to database",
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}

func UpdateCenterTagStaff(c *fiber.Ctx) error {
	var reqBody map[string]any
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request!",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request",
				Error:     err,
			},
		})
	}

	result, err := UpdateCenterStaffDB(reqBody)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Problem connecting to database",
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(result)
}
