package esystem

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetClientCurrentLoans(c *fiber.Ctx) error {
	loanDetails := new(LoanDetailsRequest)

	if err := c.BodyParser(loanDetails); err != nil {
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

	cidStr := strconv.Itoa(loanDetails.Cid)

	currentLoans, err := GetAllCurrentLoans(loanDetails.Cid, loanDetails.BrCode, loanDetails.LoanProductCode)
	if err != nil {
		logs.LOSLogs(c, LOSfeature, cidStr, "500", err.Error()+" "+loanDetails.BrCode)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "There is a problem fetching client loan details in eSystem.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, LOSfeature, cidStr, "200", "Successfully fetch client current loans. "+loanDetails.BrCode)
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successfully fetch client current loans.",
		Data:    currentLoans,
	})
}
