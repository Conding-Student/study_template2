package dashboard

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func GetDashBoardData(c *fiber.Ctx) error {
	db := database.DB

	var detailedLoanStatusCounts []map[string]any
	if err := db.Raw("SELECT * FROM loan_application.dashboarddetailedloancounts()").Scan(&detailedLoanStatusCounts).Error; err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	var loanStatusCounts map[string]any
	if err := db.Raw("SELECT * FROM loan_application.dashboardloancounts()").Scan(&loanStatusCounts).Error; err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	var designationsCount []map[string]any
	if err := db.Raw("SELECT * FROM loan_application.dashboarddesignationcount()").Scan(&designationsCount).Error; err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	var loanProducts []string
	if err := db.Raw("SELECT * FROM loan_application.dashboardactiveloanproduct()").Scan(&loanProducts).Error; err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	var loanCount []map[string]any
	if err := db.Raw("SELECT * FROM loan_application.dashboardloanamountperloantype()").Scan(&loanCount).Error; err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful",
		Data: fiber.Map{
			"detailedStatusCounts": detailedLoanStatusCounts,
			"statusCounts":         loanStatusCounts,
			"designations":         designationsCount,
			"activeLoanProd":       loanProducts,
			"loanCounts":           loanCount,
		},
	})
}
