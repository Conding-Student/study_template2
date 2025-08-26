package creditline

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func GetApprovedCreditLine(c *fiber.Ctx) error {
	creditLineReq := make(map[string]any)
	// staffID := c.Params("id")

	if err := c.BodyParser(&creditLineReq); err != nil {
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

	approvedCreditLine, err := GetCreditLineApproved(creditLineReq)
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching approved credit line",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(approvedCreditLine)
}
