package users

import (
	"chatbot/pkg/authentication"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CreatePin(c *fiber.Ctx) error {
	requestBody := make(map[string]any)

	if err := c.BodyParser(&requestBody); err != nil {
		fmt.Println("err:", err)
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	requestType := sharedfunctions.GetIntFromMap(requestBody, "reqType")

	switch requestType {
	case 0:
		staffid := sharedfunctions.GetStringFromMap(requestBody, "staffid")
		emailAddress, err := sharedfunctions.FetchUserEmail(staffid)
		if err != nil {
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   "An error occured while fetching your email address.",
					IsSuccess: false,
					Error:     err,
				},
			})
		}

		message, err := authentication.GenerateOtp(emailAddress, staffid, 2)
		if err != nil {
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   message,
					IsSuccess: false,
					Error:     err,
				},
			})
		}

		return c.JSON(response.ResponseModel{
			RetCode: "200",
			Message: "Successful!",
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: false,
				Error:     err,
			},
		})
	case 1:
		staffid := sharedfunctions.GetStringFromMap(requestBody, "staffid")
		otp := sharedfunctions.GetStringFromMap(requestBody, "otp")
		message, verified := authentication.VerifyOTP(otp, staffid, 0)
		if !verified {
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   message,
					IsSuccess: false,
					Error:     nil,
				},
			})
		}

		return c.JSON(response.ResponseModel{
			RetCode: "200",
			Message: "Successful!",
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	result, err := PinCreation(requestBody)
	if err != nil {
		fmt.Println("err:", err)
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "An error occured while processing request.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}
