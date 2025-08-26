package adminmlni

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

type ManageUser map[string]any

func UpdateMlniUser(c *fiber.Ctx) error {
	manageUser := make(ManageUser)

	if err := c.BodyParser(&manageUser); err != nil {
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

	result, isSuccess, retCodeInt, retcode, status, message, err := ManageMlniUser(manageUser)
	if err != nil {
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retcode,
			Message: status,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retcode,
		Message: status,
		Data:    result,
	})
}
