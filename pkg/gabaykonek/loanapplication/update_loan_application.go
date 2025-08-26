package loans

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/features"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UpdateLoanApplication(c *fiber.Ctx) error {
	updateLoan := make(jsonBRequestBody)
	staffID := c.Params("id")
	var module string

	if err := c.BodyParser(&updateLoan); err != nil {
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

	updateLoan["staffid"] = staffID
	customerInfo := sharedfunctions.GetMap(updateLoan, "customerInfo")
	loanHeader := sharedfunctions.GetMap(updateLoan, "loanApp")
	customerBranch := sharedfunctions.GetStringFromMap(customerInfo, "branch_code")
	customerCid := sharedfunctions.GetStringFromMap(customerInfo, "cid")
	statusCode := sharedfunctions.GetStringFromMap(loanHeader, "status_code")
	loanProductCode := sharedfunctions.GetStringFromMap(loanHeader, "loan_product_code")
	loanID := sharedfunctions.GetStringFromMap(loanHeader, "loan_id")

	var producedLoan map[string]any
	if strings.HasPrefix(statusCode, "DIS") {
		module = features.LoanDisbursementModule

		kLoan, kIsSuccess, kRetCodeInt, kRetCode, kStatus, kMessage, err := produceToKafka(updateLoan)
		if err != nil {
			logs.ErrorLogs(staffID, module, kMessage+"\nError Details: "+"\nLoan ID: "+loanID+"\nBranch Code: "+customerBranch+"\nCID: "+customerCid+"\nLoan Product Code: "+loanProductCode+"\nError: "+err.Error())
			return c.Status(kRetCodeInt).JSON(response.ResponseModel{
				RetCode: kRetCode,
				Message: kStatus,
				Data: errors.ErrorModel{
					Message:   kMessage,
					IsSuccess: kIsSuccess,
					Error:     err,
				},
			})
		}
		if !kIsSuccess {
			logs.ErrorLogs(staffID, module, kMessage+"\nError Details: "+"\nLoan ID: "+loanID+"\nBranch Code: "+customerBranch+"\nCID: "+customerCid+"\nLoan Product Code: "+loanProductCode+"\nError: Status producing loan to kafka is not successful")
			return c.Status(kRetCodeInt).JSON(response.ResponseModel{
				RetCode: kRetCode,
				Message: kStatus,
				Data: errors.ErrorModel{
					Message:   kMessage,
					IsSuccess: kIsSuccess,
					Error:     err,
				},
			})
		}

		producedLoan = kLoan
	} else if strings.HasPrefix(statusCode, "REC") {
		module = features.LoanRecommendationModule
	} else if strings.HasPrefix(statusCode, "APR") {
		module = features.LoanApprovalModule
	} else if strings.HasPrefix(statusCode, "REJ") {
		module = features.LoanRejectionModule
	}

	sharedfunctions.ConvertStringToJSONMap(updateLoan)
	isSuccess, lRetCodeInt, lRetCode, lStatus, lMessage, err := LoanUpdating(updateLoan)
	if err != nil {
		logs.ErrorLogs(staffID, module, lMessage+"\nError Details: "+"\nLoan ID: "+loanID+"\nBranch Code: "+customerBranch+"\nCID: "+customerCid+"\nLoan Product Code: "+loanProductCode+"\nError: "+err.Error())
		return c.Status(lRetCodeInt).JSON(response.ResponseModel{
			RetCode: lRetCode,
			Message: lStatus,
			Data: errors.ErrorModel{
				Message:   lMessage,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	return c.Status(lRetCodeInt).JSON(response.ResponseModel{
		RetCode: lRetCode,
		Message: lStatus,
		Data: fiber.Map{
			"producedToKafka": producedLoan,
		},
	})
}
