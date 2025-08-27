package elli

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Admin_practice_creation(c *fiber.Ctx) error {
	user_input := make(map[string]any)
	if err := c.BodyParser(&user_input); err != nil {
		fmt.Print("Trying to parse user input")
		return sharedfunctions.SendErrorResponse(c, 401, "401", status.RetCode401, "Failed to parse practice creation credentials", false, err)
	}

	// ✅ Get staffID from request
	staffID := sharedfunctions.GetStringFromMap(user_input, "staffId")

	// ✅ Call account creation logic
	isSuccess, retCodeInt, retCode, responseStatus, message, err := Creating_admin_account(user_input)
	if err != nil {
		logs.LOSLogs(c, "Account Creation Module", staffID, retCode, err.Error()+" "+staffID)
		return sharedfunctions.SendErrorResponse(c, retCodeInt, retCode, responseStatus, message, isSuccess, err)
	}

	// ✅ Log success
	logs.LOSLogs(c, "Account Creation Module", staffID, retCode, message+" "+staffID)

	// ✅ Use message + user_input as success payload
	return sharedfunctions.SendSuccessResponse(c, retCodeInt, retCode, message, user_input)
}
