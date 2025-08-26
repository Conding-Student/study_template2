package features

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func FeaturesAndVersions(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	reqBody := make(map[string]any)

	if err := c.BodyParser(&reqBody); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	var result map[string]any
	if err := db.Raw("SELECT * FROM  public.featureversions($1)", reqBody).Scan(&result).Error; err != nil {
		logs.LOSLogs(c, "AppOpenning", id, "500", err.Error())
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to fetch features and versions.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "featureversions")

	logs.LOSLogs(c, "AppOpenning", id, "200", "Successfully fetch initial data.")
	return c.JSON(result)
}
