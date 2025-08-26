package authentication

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ValidateSession(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		fmt.Println("Missing or invalid token: ", authHeader)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		fmt.Println("Invalid token format: ", authHeader)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	isSuccess, retCodeInt, retCode, tstatus, tmessage, err := sharedfunctions.ValidateToken(tokenString)
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

	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: tmessage,
		Data:    nil,
	})
}
