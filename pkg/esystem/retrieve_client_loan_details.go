package esystem

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"log"

	"github.com/gofiber/fiber/v2"
)

type LoanDetailsRequest struct {
	Cid             int    `gorm:"not null"`
	LoanProductCode int    `gorm:"not null"`
	BrCode          string `gorm:"not null"`
}

func GetClientLoanDetails(c *fiber.Ctx) error {
	loanReqDetails := make(map[string]any)

	if err := c.BodyParser(&loanReqDetails); err != nil {
		log.Println(err)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request body.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	cidStr := sharedfunctions.GetStringFromMap(loanReqDetails, "Cid")
	brcode := sharedfunctions.GetStringFromMap(loanReqDetails, "BrCode")

	previousLoan, err := GetClientPreviousLoan(loanReqDetails)
	if err != nil {
		logs.LOSLogs(c, LOSfeature, cidStr, "500", err.Error()+" "+brcode)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching client previous loan.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, LOSfeature, cidStr, "200", "Client Loan Information retrieved successfully! "+brcode)
	return c.JSON(previousLoan)
}
