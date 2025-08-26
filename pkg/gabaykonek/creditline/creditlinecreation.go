package creditline

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NewCreditLine struct {
	CreditLineFields map[string]any `json:"creditLineFields"`
}

func CreditLineCreation(c *fiber.Ctx) error {
	newCreditLine := new(NewCreditLine)
	operation := c.Get("operation")
	staffID := c.Params("id")

	if err := c.BodyParser(newCreditLine); err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	personalInfoFields := sharedfunctions.GetMap(newCreditLine.CreditLineFields, "personalInfoFields")
	brcode := sharedfunctions.GetStringFromMap(personalInfoFields, "brcode")
	unit := sharedfunctions.GetIntFromMap(personalInfoFields, "unit")
	centercode := sharedfunctions.GetStringFromMap(personalInfoFields, "centercode")
	cid := sharedfunctions.GetIntFromMap(personalInfoFields, "cid")
	ref := sharedfunctions.GetStringFromMap(personalInfoFields, "ref")

	isSuccess, retCodeInt, retCode, status, message, err := CreateCreditLine(operation, newCreditLine)
	if err != nil {
		logs.ErrorLogs(staffID, module, message+"\nError: "+err.Error()+"\nnCreation Details: "+"\nReferance :"+ref+"\nBranch :"+brcode+"\nUnit: "+strconv.Itoa(unit)+"\nCenter: "+centercode+"\nCID: "+strconv.Itoa(cid))
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: status,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.CardIncAuditTrail(staffID, module, message+"\nCreation Details: "+"\nReferance :"+ref+"\nBranch :"+brcode+"\nUnit: "+strconv.Itoa(unit)+"\nCenter: "+centercode+"\nCID: "+strconv.Itoa(cid))
	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: status,
		Data: errors.ErrorModel{
			IsSuccess: isSuccess,
			Message:   message,
			Error:     nil,
		},
	})
}
