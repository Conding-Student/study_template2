package users

import (
	"chatbot/pkg/authentication"
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

type UpdateDeviceRequest map[string]any

var module = "Update Device Module"

func UpdateDevice(c *fiber.Ctx) error {
	staffID := c.Params("id")
	otp := c.Get("otp")

	updateDevice := make(UpdateDeviceRequest)

	if err := c.BodyParser(&updateDevice); err != nil {
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

	requestType := sharedfunctions.GetIntFromMap(updateDevice, "RequestType")

	if requestType == 0 {

		email, err := sharedfunctions.FetchUserEmail(staffID)
		if err != nil {
			logs.LOSLogs(c, module, staffID, "401", email+" "+staffID)
			return c.Status(401).JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   email,
					IsSuccess: false,
					Error:     nil,
				},
			})
		}

		message, err := authentication.GenerateOtp(email, staffID, 1)
		if err != nil {
			logs.LOSLogs(c, module, staffID, "401", "We encounter an error while sending a One Time Password (OTP) in your email. Please try again later. "+staffID)
			return c.Status(401).JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   message,
					IsSuccess: false,
					Error:     err,
				},
			})
		}

		logs.LOSLogs(c, module, staffID, "200", "OTP sent to user email for device updating purposes.")
		return c.Status(200).JSON(response.ResponseModel{
			RetCode: "200",
			Message: "Successful!",
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: true,
				Error:     nil,
			},
		})

	}

	message, verified := authentication.VerifyOTP(otp, staffID, 0)
	if !verified {
		logs.LOSLogs(c, module, staffID, "200", "OTP was invalid.")
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

	isSuccess, retCodeInt, retCode, status, message, err := sharedfunctions.UpdateDevice(updateDevice)
	if err != nil {
		logs.LOSLogs(c, module, staffID, retCode, message)
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

	logs.LOSLogs(c, module, staffID, retCode, message)
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
