package loans

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func ViewMidasDetails(c *fiber.Ctx) error {
	LOSFeature := "LOS - Retrieve Midas Details"

	var results map[string]any
	if err := database.DB.Raw("SELECT * FROM gabaykonekfunc.getmidas()").Scan(&results).Error; err != nil {
		database.DB.Rollback()
		logs.LOSLogs(c, LOSFeature, "", "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(results)
	midas := sharedfunctions.GetList(results, "getmidas")

	logs.LOSLogs(c, LOSFeature, "", "200", "Midas data retrieved successfully")
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Midas data retrieved successfully",
		Data:    midas,
	})
}
