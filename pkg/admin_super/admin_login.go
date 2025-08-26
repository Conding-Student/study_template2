package administrator

import (
	"chatbot/pkg/logs"
	//"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"

	//"chatbot/pkg/models/errors"

	"github.com/gofiber/fiber/v2"
)

var adminLoginModule = "Admin Login Module"

// func AdminLogin(c *fiber.Ctx) error {
// 	staffID := c.Params("id")

// 	loginCreds := make(map[string]any)

// 	if err := c.BodyParser(&loginCreds); err != nil {
// 		fmt.Println("Failed to parse login credentials", err)
// 		return c.Status(401).JSON(response.ResponseModel{
// 			RetCode: "401",
// 			Message: status.RetCode401,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to parse login credentials",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	userData, isSuccess, retCodeInt, retCode, responseStatus, loginMessage, err := sharedfunctions.AccountLoginAdmin(staffID, loginCreds)
// 	if err != nil {
// 		logs.LOSLogs(c, adminLoginModule, staffID, retCode, loginMessage)
// 		return c.Status(retCodeInt).JSON(response.ResponseModel{
// 			RetCode: retCode,
// 			Message: responseStatus,
// 			Data: errors.ErrorModel{
// 				Message:   loginMessage,
// 				IsSuccess: isSuccess,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	logs.LOSLogs(c, adminLoginModule, staffID, retCode, loginMessage)
// 	return c.Status(retCodeInt).JSON(response.ResponseModel{
// 		RetCode: retCode,
// 		Message: loginMessage,
// 		Data:    userData,
// 	})

// }

func AdminLogin(c *fiber.Ctx) error {
	staffID := c.Params("id")

	loginCreds := make(map[string]any)

	if err := c.BodyParser(&loginCreds); err != nil {
		fmt.Println("Failed to parse login credentials", err)
		//parameters for send error response: fiber, status code int, retcode string, message string, isSuccess bool, err error
		return sharedfunctions.SendErrorResponse(c, 401, "401", status.RetCode401, "Failed to parse login credentials", false, err)
	}

	userData, isSuccess, retCodeInt, retCode, responseStatus, loginMessage, err := sharedfunctions.AccountLoginAdmin(staffID, loginCreds)
	if err != nil {
		logs.LOSLogs(c, adminLoginModule, staffID, retCode, loginMessage)
		return sharedfunctions.SendErrorResponse(c, retCodeInt, retCode, responseStatus, loginMessage, isSuccess, err)
	}

	logs.LOSLogs(c, adminLoginModule, staffID, retCode, loginMessage)
	fmt.Println("updated login")
	//parameters for send success response: fiber, red code int, retcode string, message string, data interface{}
	return sharedfunctions.SendSuccessResponse(c, retCodeInt, retCode, loginMessage, userData)

}
