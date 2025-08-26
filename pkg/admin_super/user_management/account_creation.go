package usermanagement

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

func AccountCreationAdmin(c *fiber.Ctx) error {
	staffInfo := make(map[string]any)

	if err := c.BodyParser(&staffInfo); err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse change password credentials",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	staffID := sharedfunctions.GetStringFromMap(staffInfo, "staffId")

	isSuccess, retCodeInt, retCode, responseStatus, message, err := sharedfunctions.AccountCreationAdmin(staffInfo)
	if err != nil {
		logs.LOSLogs(c, "Account Creation Module", staffID, retCode, err.Error()+" "+staffID)
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: responseStatus,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, "Account Creation Module", staffID, retCode, message+" "+staffID)
	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: responseStatus,
		Data: fiber.Map{
			"isSuccess": isSuccess,
			"message":   message,
			"error":     nil,
		},
	})
}
