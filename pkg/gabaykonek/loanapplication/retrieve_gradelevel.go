package loans

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func GetGradeLevel(c *fiber.Ctx) error {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getbengradelvl()").Scan(&result).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to fetch grade level.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	gradeLevel := sharedfunctions.GetList(result, "response")

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Grade level fetch successfully!",
		Data:    gradeLevel,
	})
}
