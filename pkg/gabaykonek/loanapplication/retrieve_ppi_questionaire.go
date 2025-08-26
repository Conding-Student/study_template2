package loans

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func GetPPIQuestionaire(c *fiber.Ctx) error {
	LOSFeature := "LOS - PPI Questionares"

	var results map[string]any
	if err := database.DB.Raw("SELECT * FROM gabaykonekfunc.getppiquestionaire()").Scan(&results).Error; err != nil {
		database.DB.Rollback()
		logs.LOSLogs(c, LOSFeature, "PPI Questionares retrieving failed", "500", err.Error())
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

	sharedfunctions.ConvertStringToJSONMap(results)
	questionaire := sharedfunctions.GetList(results, "getppiquestionaire")

	logs.LOSLogs(c, LOSFeature, "", "200", "PPI Questionares retrieved successfully")
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "PPI Questionares",
		Data:    questionaire,
	})
}
