package users

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ProfilePictureCredential struct {
	StaffID string `json:"staffid"`
	Picture string `json:"picture"`
}

func UpdateUserProfilePicture(c *fiber.Ctx) error {
	staffID := c.Params("id")

	var requestBody ProfilePictureCredential
	if err := c.BodyParser(&requestBody); err != nil {
		fmt.Println("Failed to parse request", err)
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

	isSuccess, message, err := UpdateProfilePic(requestBody.StaffID, requestBody.Picture)
	if err != nil {
		logs.LOSLogs(c, "UpdateProfilePicture", staffID, "500", message)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, "UpdateProfilePicture", staffID, "200", "Profile picture updated successfully.")
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: message,
		Data:    nil,
	})
}
