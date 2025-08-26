package users

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	staffid := c.Params("id")
	params := make(map[string]any)

	if err := c.BodyParser(&params); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request.",
				Error:     err,
			},
		})
	}

	params["paramStaffid"] = staffid
	result, err := sharedfunctions.AccountLogout(params)
	if err != nil {
		logs.LOSLogs(c, "Logout", staffid, "500", err.Error())
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "An error occured while connecting to database.",
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}
