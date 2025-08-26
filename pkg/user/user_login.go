package users

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	staffID := c.Params("id")
	deviceID := c.Params("deviceid")

	loginCreds := make(map[string]any)

	if err := c.BodyParser(&loginCreds); err != nil {
		fmt.Println("Failed to parse login credentials", err)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse login credentials",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	userData, isSuccess, retCodeInt, retCode, responseStatus, loginMessage, err := sharedfunctions.AccountLoginV2(loginCreds, staffID, deviceID)
	if err != nil {
		logs.LOSLogs(c, "Login", staffID, retCode, loginMessage)
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: responseStatus,
			Data: errors.ErrorModel{
				Message:   loginMessage,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, "Login", staffID, retCode, loginMessage)
	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: responseStatus,
		Data:    userData,
	})
}
