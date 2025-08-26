package loancalc

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func RBILoanCalculator(c *fiber.Ctx) error {
	id := c.Params("id")
	computeLoan := make(map[string]any)
	LOSFeature := "RBI Loan calculator"

	if err := c.BodyParser(&computeLoan); err != nil {
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

	computeLoan["id"] = id
	result, err := RBICalcForBanks(computeLoan)
	if err != nil {
		logs.LOSLogs(c, LOSFeature, id, "500", "An error occure while computing loan."+" "+err.Error())
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occure while computing loan.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, LOSFeature, id, "200", "Loan computation successful!")
	return c.JSON(result)
}
