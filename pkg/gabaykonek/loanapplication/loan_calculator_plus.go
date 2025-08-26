package loans

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Compute map[string]any

type LoanCalculatorPlusResponse struct {
	LoanAmount      string           `json:"loanAmount"`
	EIR             string           `json:"effectiveInterestRate"`
	ContractualRate string           `json:"contractualRate"`
	Interest        string           `json:"loanInterest"`
	InterestRate    string           `json:"loanInterestRate"`
	LoanOutstanding string           `json:"loanOutstanding"`
	NumberOfMonths  string           `json:"numberOfMonths"`
	NumberOfWeeks   string           `json:"numberOfWeeks"`
	LRF             string           `json:"lrf"`
	LoanProceeds    string           `json:"loanProceeds"`
	AmountInWords   string           `json:"amountInWords"`
	WeeklyDue       string           `json:"weeklyDue"`
	ReleaseDate     string           `json:"dateRelease"`
	StartDate       string           `json:"firstPayment"`
	MaturityDate    string           `json:"maturityDate"`
	Amortization    []map[string]any `json:"amortization"`
}

func LoanCalculatorPlus(c *fiber.Ctx) error {
	computeLoan := make(Compute)
	LOSFeature := "LOS - kPlus Loan Calculator"

	if err := c.BodyParser(&computeLoan); err != nil {
		fmt.Println("Failed to to parse request", err)
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

	instiCode := sharedfunctions.GetStringFromMap(computeLoan, "instiCode")

	result, isSuccess, retCodeInt, retCode, status, message, err := LoanCalculator(computeLoan)
	if err != nil {
		logs.LOSLogs(c, LOSFeature, instiCode, "500", message+" "+err.Error())
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: status,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, LOSFeature, instiCode, "200", "Loan computation successful!")
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Loan computation successful!",
		Data:    result,
	})
}
