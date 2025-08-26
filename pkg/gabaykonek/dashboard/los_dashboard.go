package gabaykonekdashboard

import (
	"chatbot/pkg/gabaykonek/reports"
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/features"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

type LoansAndTotal struct {
	StaffID   string
	StartDate string
	EndDate   string
}

func LOSDashboard(c *fiber.Ctx) error {
	db := database.DB
	staffID := c.Params("id")
	module := features.DashboardModule

	query := "SELECT * FROM gabaykonekfunc.getloanstatuscount($1)"

	var dashboardData map[string]any
	if err := db.Raw(query, staffID).Scan(&dashboardData).Error; err != nil {
		logs.ErrorLogs(staffID, module, "The user fails to fetched dashboard data\n Error Details: User not permitted to view dashboard data or not authorized to view the feature.")
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "You are not permitted to view the data.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(dashboardData)

	data := sharedfunctions.GetMap(dashboardData, "getloanstatuscount")

	logs.CardIncAuditTrail(staffID, module, "User successfully fetched dashboard data")
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successfully fetch dashboard data",
		Data:    data,
	})
}

func LoansAndTotals(c *fiber.Ctx) error {
	module := features.DashboardModule
	loansAndTotas := new(LoansAndTotal)

	if err := c.BodyParser(loansAndTotas); err != nil {
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

	staffID := loansAndTotas.StaffID

	dateRange, err := reports.GetDateRange(loansAndTotas.StartDate, loansAndTotas.EndDate)
	if err != nil {
		logs.ErrorLogs(staffID, module, "Failed to generate Summary of Loan Releases. "+dateRange+"\nError Details: "+err.Error())
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

	summary, err := GetReleasedLoanAndTotal(staffID, loansAndTotas.StartDate, loansAndTotas.EndDate)
	if err != nil {
		logs.ErrorLogs(staffID, module, "The user fails to fetched released loans and total amount from "+loansAndTotas.StartDate+" to "+loansAndTotas.EndDate+"\nError Details: "+err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to fetch loan released and amount.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.CardIncAuditTrail(staffID, module, "User successfully fetched Released Loan From "+loansAndTotas.StartDate+" to "+loansAndTotas.EndDate)
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successfully fetch status count.",
		Data: fiber.Map{
			"dateRange":      "Successful Releases " + dateRange,
			"loansAndCounts": summary,
		},
	})
}
