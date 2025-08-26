package users

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

type ChangePasswordCredentials struct {
	StaffID         string
	CurrentPassword string
	NewPassword     string
}

func ChangePassword(c *fiber.Ctx) error {
	staffID := c.Params("id")

	changePasswordCredentials := new(ChangePasswordCredentials)
	if err := c.BodyParser(changePasswordCredentials); err != nil {
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

	isSuccess, retCodeInt, retCode, status, message, err := sharedfunctions.UpdatePassword(changePasswordCredentials.StaffID, changePasswordCredentials.CurrentPassword, changePasswordCredentials.NewPassword)
	if err != nil {
		logs.LOSLogs(c, "ChangePassword", staffID, retCode, message)
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: status,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     nil,
			},
		})
	}

	logs.LOSLogs(c, "ChangePassword", staffID, retCode, message)
	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: message,
		Data: errors.ErrorModel{
			Message:   message,
			IsSuccess: isSuccess,
			Error:     nil,
		},
	})
}
