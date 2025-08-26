package loans

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetCoBorrowerRelationships(c *fiber.Ctx) error {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getrelationships()").Scan(&result).Error; err != nil {
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
	relationships := sharedfunctions.GetMap(result, "response")
	relationshipList := sharedfunctions.GetListString(relationships, "relationships")
	fmt.Println(relationships)

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Relationships fetch successfully!",
		Data:    relationshipList,
	})
}
