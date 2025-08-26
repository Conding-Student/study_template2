package empc

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetEMPCAmortization(c *fiber.Ctx) error {
	empcLoanAmortRequest := make(map[string]any)

	if err := c.BodyParser(&empcLoanAmortRequest); err != nil {
		fmt.Println("retCode 401")
		fmt.Println("Invalid Request")
		fmt.Println("Failed to parse request", err.Error())
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

	staffId := sharedfunctions.GetStringFromMap(empcLoanAmortRequest, "staffid")

	loanAmortizationResponseData, err := GetEMPCLoan(empcLoanAmortRequest)
	if err != nil {
		logs.LOSLogs(c, EMPCFeature, staffId, "500", err.Error())
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching loan amortization. Please try again later.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(loanAmortizationResponseData)
}
