package loans

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/features"
	"chatbot/pkg/models/response"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

func GetLoanApplications(c *fiber.Ctx) error {
	module := features.LoanRetrievalModule
	staffID := c.Params("id")

	loanApplication, isSuccess, retCodeInt, retCode, gStatus, message, err := GetAllLoans(staffID)
	if err != nil {
		logs.ErrorLogs(staffID, module, message+"\nDetails: "+err.Error())
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: gStatus,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	loans := sharedfunctions.GetListAny(loanApplication, "data")

	logs.CardIncAuditTrail(staffID, module, "The user successfully retrieve list of loans.\nDetails: "+message)
	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: gStatus,
		Data:    loans,
	})
}
