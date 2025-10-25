package loans

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

type LoanProductResponse struct {
	LoanProducts string `json:"loan_products"`
}

func LoanProducts(c *fiber.Ctx) error {
	var result map[string]any
	if err := database.DB.Raw("SELECT * FROM gabaykonekfunc.loan_products()").Scan(&result).Error; err != nil {
		log.Println(err)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	loans := sharedfunctions.GetMap(result, "loan_products")

	return c.JSON(loans)
}

func LoanProductListAndDetails(c *fiber.Ctx) error {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getactiveloans()").Find(&result).Error; err != nil {
		log.Println(err)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch loan products",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	activeLoans := sharedfunctions.GetList(result, "getactiveloans")

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    activeLoans,
	})
}
