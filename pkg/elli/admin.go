package elli

import (
	//"chatbot/pkg/logs"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	//"chatbot/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func admin_practice_creation(c *fiber.Ctx) error {
	user_input := make(map[string]any)
	if err := c.BodyParser(&user_input); err != nil {
		fmt.Print("Trying to parse user input")
		//parameters for send error response: fiber, status code int, retcode string, message string, isSuccess bool, err error
		return sharedfunctions.SendErrorResponse(c, 401, "401", status.RetCode401, "Failed to parse practice creation credentials", false, err)
	}

	// staffInfo := sharedfunctions.GetStringFromMap(user_input, "staffId")
	// isSuccess, retCodeInt, retCode, responseStatus, message, err := sharedfunctions.AccountCreationAdmin(staffInfo)

}
