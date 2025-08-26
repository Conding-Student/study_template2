package eloading

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func EloadLoadRequest(c *fiber.Ctx) error {
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

	requestType := sharedfunctions.GetIntFromMap(requestBody, "requestType")
	if requestType == 3 {
		validateResp, err := sharedfunctions.PinValidation(requestBody)
		if err != nil {
			fmt.Println("err:", err)
			return c.JSON(response.ResponseModel{
				RetCode: "500",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   "An error occured while connecting to database.",
					IsSuccess: false,
					Error:     err,
				},
			})
		}

		if sharedfunctions.GetStringFromMap(validateResp, "retCode") != "200" {
			return c.JSON(validateResp)
		}
	}

	result, err := serviceRequest(requestBody)
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
