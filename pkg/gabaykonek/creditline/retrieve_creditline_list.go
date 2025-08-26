package creditline

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

type GetCreditLineReq struct {
	StaffID string
}

func GetCreditLineList(c *fiber.Ctx) error {
	getCreditLine := new(GetCreditLineReq)
	staffID := c.Params("id")

	if err := c.BodyParser(&getCreditLine); err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	creditLine, err := GetCreditLine(getCreditLine.StaffID)
	if err != nil {
		logs.ErrorLogs(staffID, module, "An error occured while fetching credit line list."+"\nError: "+err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching credit line list.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.CardIncAuditTrail(staffID, module, "User successfully fetch credit line list.")
	return c.JSON(creditLine)
}
