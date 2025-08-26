package empc

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

type EMPCLoanProduct struct {
	LoanID       string `json:"loan_id"`
	LoanProducts string `json:"loan_products"`
	Description  string `json:"description"`
}

func GetEmpcLoanProducts(c *fiber.Ctx) error {
	db := database.DB

	staffId := c.Params("id")
	// Fetch all rows from the database table
	rows, err := db.Raw("SELECT * FROM empc.empc_loans").Rows()
	if err != nil {
		logs.LOSLogs(c, EMPCFeature, staffId, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}
	defer rows.Close()

	// Create a slice to hold the results
	var loanProducts []EMPCLoanProduct

	// Iterate over the result set
	for rows.Next() {
		var loanProduct EMPCLoanProduct
		if err := rows.Scan(&loanProduct.LoanID, &loanProduct.LoanProducts, &loanProduct.Description); err != nil {
			logs.LOSLogs(c, EMPCFeature, staffId, "500", err.Error())
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Internal server error",
				Data: errors.ErrorModel{
					Message:   "Failed to scan rows",
					IsSuccess: false,
					Error:     err,
				},
			})
		}
		// Append the current row to the slice
		loanProducts = append(loanProducts, loanProduct)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		logs.LOSLogs(c, EMPCFeature, staffId, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Error iterating rows",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    loanProducts,
	})
}
