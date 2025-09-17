package loanreversal

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func ViewLoans(c *fiber.Ctx) error {
	reqBody := make(map[string]any)

	if err := c.BodyParser(&reqBody); err != nil {
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

	// Call helper function
	result, err := loan_viewing(reqBody)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}

// func ViewLoans(c *fiber.Ctx) error {
// 	db := database.DB
// 	reqBody := make(map[string]any)

// 	if err := c.BodyParser(&reqBody); err != nil {
// 		return c.JSON(response.ResponseModel{
// 			RetCode: "401",
// 			Message: status.RetCode401,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to parse request",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	var result map[string]any
// 	if err := db.Raw("SELECT gabaykonekfunc.getlistofdisloan($1)", reqBody).Scan(&result).Error; err != nil {
// 		return c.JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: status.RetCode500,
// 			Data: errors.ErrorModel{
// 				Message:   "An error occured while connecting to database",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	sharedfunctions.ConvertStringToJSONMap(result)
// 	loans := sharedfunctions.GetMap(result, "getlistofdisloan")
// 	return c.JSON(loans)
// }
