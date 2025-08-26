package reports

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/features"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetPendingLoanSummary(c *fiber.Ctx) error {
	summaryRequest := new(ReportsRequestBody)

	if err := c.BodyParser(summaryRequest); err != nil {
		log.Println(err)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	staffID := summaryRequest.Staffid
	module := features.ReportsModule
	dateRange, err := GetDateRange(summaryRequest.StartDate, summaryRequest.EndDate)
	if err != nil {
		logs.ErrorLogs(staffID, module, "Failed to generate Pending Loan Application. "+dateRange+"\nError Details: "+err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   err.Error(),
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	desigInt, err := sharedfunctions.DesignationToInt(summaryRequest.Designation)
	if err != nil {
		logs.ErrorLogs(staffID, module, "Failed to generate Pending Loan Application. "+dateRange+"\nError Details: "+err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   err.Error(),
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	summaryOfLoans, err := GetPendingSummary(staffID, desigInt, summaryRequest.StartDate, summaryRequest.EndDate)
	if err != nil {
		logs.ErrorLogs(staffID, module, "Failed to generate Pending Loan Application. "+dateRange+"\nError Details: "+err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to fetch summary of reports.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	branchName, err := GetBranchName(summaryRequest.Staffid, desigInt)
	if err != nil {
		logs.ErrorLogs(staffID, module, "Failed to generate Pending Loan Application. "+dateRange+"\nError Details: "+err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to fetch summary of reports.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	dateTime, err := sharedfunctions.LocalTime()
	if err != nil {
		logs.ErrorLogs(staffID, module, "Failed to generate Pending Loan Application. "+dateRange+"\nError Details: "+err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to generate summary of reports.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	fullname, err := GetGeneratedBy(summaryRequest.Staffid)
	if err != nil {
		logs.ErrorLogs(staffID, module, "Failed to identify staff generating Pending Loan Application. "+dateRange+"\nError Details: "+err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to identify staff generating summary of reports.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.CardIncAuditTrail(staffID, module, "The user successfully fetched Summary of Pending Loan Application "+dateRange)
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successfully fetch summary of reports",
		Data: fiber.Map{
			"reportstitle":  "Summary of Pending Loan Application",
			"dategenerated": dateTime,
			"generatedby":   fullname,
			"branch":        branchName,
			"daterange":     dateRange,
			"reportsdata":   summaryOfLoans,
		},
	})
}
