package loans

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type LoanStatus struct {
	LoanID string
}

// For CA-GABAY and kPlus
func GetLoanStatus(c *fiber.Ctx) error {
	db := database.DB
	LOSFeature := "LOS - Retrieve Loan Status"

	loanStatusRequest := new(LoanStatus)
	if err := c.BodyParser(loanStatusRequest); err != nil {
		log.Println(err)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Unable to load loan status.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getloanstatus($1)", loanStatusRequest.LoanID).Scan(&result).Error; err != nil {
		log.Println(err)
		logs.LOSLogs(c, LOSFeature, loanStatusRequest.LoanID, "500", "Unable to load loan status. "+loanStatusRequest.LoanID+" "+err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching loan status.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	responseData := sharedfunctions.GetMap(result, "response")

	isSuccess := sharedfunctions.GetBoolFromMap(responseData, "issucces")
	lStatus := sharedfunctions.GetStringFromMap(responseData, "status")
	retCode := sharedfunctions.GetStringFromMap(responseData, "retcode")
	retCodeInt := sharedfunctions.GetIntFromMap(responseData, "retcode")
	message := sharedfunctions.GetStringFromMap(responseData, "message")
	loanStatusList := sharedfunctions.GetMap(responseData, "data")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nLoan status fetching successful: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	if !isSuccess {
		logs.LOSLogs(c, LOSFeature, loanStatusRequest.LoanID, retCode, message+"\n"+loanStatusRequest.LoanID)
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: lStatus,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     fmt.Errorf(message),
			},
		})
	}

	logs.LOSLogs(c, LOSFeature, loanStatusRequest.LoanID, retCode, message+"\n"+loanStatusRequest.LoanID)
	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: lStatus,
		Data:    loanStatusList,
	})
}

// for CA-GABAY and ADMIN
func GetLoansPerClient(c *fiber.Ctx) error {
	db := database.DB
	LOSFeature := "LOS - Retrieve Loan Status"
	id := c.Params("id")

	params := make(map[string]any)
	if err := c.BodyParser(&params); err != nil {
		log.Println(err)
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Unable to load loan status.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	params["staffid"] = id
	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getloansperclient($1)", params).Scan(&result).Error; err != nil {
		log.Println(err)
		logs.LOSLogs(c, LOSFeature, id, "500", "Unable to load loan status. "+err.Error())
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching loan status.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "response")

	return c.JSON(result)
}
