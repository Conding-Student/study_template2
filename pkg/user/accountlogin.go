package users

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/model"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	jwtToken "chatbot/pkg/utils/go-utils/fiber"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AccountLogin(c *fiber.Ctx) error {
	login := new(model.LoginCredentials)

	if err := c.BodyParser(&login); err != nil {
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

	// isSuccess, message, errDetail, err := sharedfunctions.AccountLogins(login)
	// if err != nil {
	// 	return c.Status(401).JSON(response.ResponseModel{
	// 		RetCode: "401",
	// 		Message: status.RetCode401,
	// 		Data: errors.ErrorModel{
	// 			Message:   message,
	// 			IsSuccess: isSuccess,
	// 			Error:     err,
	// 		},
	// 	})
	// }

	password := login.Password
	isSuccess, staffid, username, deviceModel, deviceId, errDetail, rolename, err := sharedfunctions.AccountDetailsForSoteria(login.StaffId)
	if err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   errDetail,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	fmt.Println(rolename)
	accessToken, err := jwtToken.GenerateJWTSignedString(fiber.Map{
		"staffID":  staffid,
		"username": username,
		"rolename": rolename,
	})

	if err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to generate token",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	isSuccess, retCodeInt, retCode, tstatus, tmessage, err := sharedfunctions.SaveTokenToDB(staffid, accessToken)
	if err != nil {
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: tstatus,
			Data: errors.ErrorModel{
				Message:   tmessage,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	message, token, err := sharedfunctions.SoteriaLogin(username, password, deviceId, deviceModel)
	if err != nil {
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

	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: tstatus,
		Data: fiber.Map{
			"soteriaToken": token,
			"cagabayToken": accessToken,
			"isSuccess":    isSuccess,
			"staffid":      staffid,
			"username":     username,
			"deviceModel":  deviceModel,
			"deviceId":     deviceId,
			"error":        errDetail,
		},
	})
}
