package features

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func SecondaryFeatures(c *fiber.Ctx) error {
	params := make(map[string]any)

	if err := c.BodyParser(&params); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse change password credentials",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	secFeatures, err := GetSecFeatures()
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching secondary features to database.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(secFeatures)
}
