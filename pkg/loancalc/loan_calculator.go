package loancalc

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Compute map[string]any

type LoanResponse struct {
	LoanAmount       string           `json:"loanAmount"`
	Contractual      string           `json:"contractual"`
	EIR              string           `json:"effectiveInterestRate"`
	Interest         string           `json:"loanInterest"`
	NumberOfMonths   string           `json:"numberOfMonths"`
	NumberOfWeeks    string           `json:"numberOfWeeks"`
	LRF              string           `json:"lrf"`
	LoanBalance      string           `json:"loanBalance"`
	DocumentaryStamp string           `json:"documentaryStamp"`
	TotalDeduction   string           `json:"totalDeduction"`
	LoanProceeds     string           `json:"loanProceeds"`
	WeeklyDue        string           `json:"weeklyDue"`
	ReleaseDate      string           `json:"dateRelease"`
	StartDate        string           `json:"firstPayment"`
	MaturityDate     string           `json:"maturityDate"`
	AmountInWords    string           `json:"amountInWords"`
	Amortization     []map[string]any `json:"amortization"`
}

func BankLoanCalculator(c *fiber.Ctx) error {
	staffID := c.Params("id")
	computeLoan := make(Compute)
	LOSFeature := "Loan calculator for all"

	if err := c.BodyParser(&computeLoan); err != nil {
		fmt.Println("Failed to to parse request", err)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	computeLoan["staffID"] = staffID
	result, isSuccess, retCodeInt, retCode, status, message, err := LoanCalcForBanks(computeLoan)
	if err != nil {
		logs.LOSLogs(c, LOSFeature, staffID, "500", message+" "+err.Error())
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

	logs.LOSLogs(c, LOSFeature, staffID, retCode, "Loan computation successful!")
	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: status,
		Data:    result,
	})
}
