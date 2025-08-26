package administrator

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func GetInstiAndClientCount(c *fiber.Ctx) error {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM public.getsuperaddashboard()").Scan(&result).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "An error occured while fetching dashboard data",
				Error:     nil,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getsuperaddashboard")
	return c.JSON(result)
}
