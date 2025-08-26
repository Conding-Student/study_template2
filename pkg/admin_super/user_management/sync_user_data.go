package usermanagement

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SyncRequest struct {
	StaffID string
}

func SyncUserData(c *fiber.Ctx) error {
	instiCode := c.Get("instiCode")
	request := new(SyncRequest)

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

	fmt.Println("Insti Code: ", instiCode)

	isSuccess, retCodeInt, retCode, responseStatus, message, err := sharedfunctions.SyncUserData(request.StaffID)
	if err != nil {
		fmt.Println("Syncing error: ", err)
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

	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: responseStatus,
		Data: errors.ErrorModel{
			Message:   message,
			IsSuccess: isSuccess,
			Error:     nil,
		},
	})
}
