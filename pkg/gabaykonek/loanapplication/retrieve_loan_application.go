package loans

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/features"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

func GetLoanApplications(c *fiber.Ctx) error {
	module := features.LoanRetrievalModule
	staffID := c.Params("id")

	loanApplication, err := GetAllLoans(staffID)
	message := sharedfunctions.GetStringFromMap(loanApplication, "message")
	if err != nil {
		logs.ErrorLogs(staffID, module, message+"\nDetails: "+err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "failed to fetch all loans",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	//loans := sharedfunctions.GetListAny(loanApplication, "data")
	logs.CardIncAuditTrail(staffID, module, "The user successfully retrieve list of loans.\nDetails: "+message)
	return c.JSON(loanApplication)
}
