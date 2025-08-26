package admincardinc

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetLoanReleased(c *fiber.Ctx) error {
	db := database.DB
	dates := make(map[string]any)
	staffid := c.Params("id")

	if err := c.BodyParser(&dates); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request",
				Error:     nil,
			},
		})
	}

	dates["staffid"] = staffid
	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getloandisbursed($1)", dates).Scan(&result).Error; err != nil {
		fmt.Println(err)
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "An error occured while fetching list of disbursed loans.",
				Error:     nil,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getloandisbursed")
	return c.JSON(result)
}
