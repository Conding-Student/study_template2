package usermanagement

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	RequestData map[string]any
}

func UpdateUsers(c *fiber.Ctx) error {
	adminAccess := c.Get("adminAccess")

	request := new(Request)

	if err := c.BodyParser(&request); err != nil {
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

	result, err := sharedfunctions.UpdateUser(adminAccess, request.RequestData)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Error updating user due to a problem connecting to database!",
				IsSuccess: false,
				Error:     err,
			},
		})
	}
	return c.JSON(result)
}
