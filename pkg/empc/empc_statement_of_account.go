package empc

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func ViewStatementOfAccount(c *fiber.Ctx) error {
	params := make(map[string]any)
	if err := c.BodyParser(&params); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	staffId := c.Params("id")
	statementOfAccount, err := GetSOA(params)
	if err != nil {
		logs.LOSLogs(c, EMPCFeature, staffId, "500", err.Error())
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching statement of account.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(statementOfAccount)
}
