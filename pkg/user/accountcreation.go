package users

import (
	"chatbot/pkg/authentication"
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type StaffInformation map[string]any

func AccountCreation(c *fiber.Ctx) error {
	staffID := c.Params("id")
	// deviceID := c.Params("deviceid")
	otp := c.Get("otp")

	staffInfo := make(StaffInformation)

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

	// fmt.Println("First Stage: ", staffInfo)

	fmt.Println("Otp from header: ", otp)
	fmt.Println("staff ID: ", staffID)
	message, verified := authentication.VerifyOTPForAccountCreation(otp, staffID, 1)
	if !verified {
		logs.LOSLogs(c, "Account Creation Module", staffID, "401", "Invalid OTP "+staffID)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	isSuccess, retCodeInt, retCode, responseStatus, message, err := sharedfunctions.AccountCreationMobile(staffInfo)
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
