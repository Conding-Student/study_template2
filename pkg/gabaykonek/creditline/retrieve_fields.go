package creditline

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func GetCreditLineFields(c *fiber.Ctx) error {

	fields, err := GetCredLineFields()
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to fetch Business Assets fields",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successfully fetch fields",
		Data: fiber.Map{
			"creditLineFields": fields,
		},
	})
}

func GetCreditLineProperties(c *fiber.Ctx) error {

	texts, fonts, retcodeInt, retcode, status, message, err := GetCredLineProperties()
	if err != nil {
		return c.Status(retcodeInt).JSON(response.ResponseModel{
			RetCode: retcode,
			Message: status,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(retcodeInt).JSON(response.ResponseModel{
		RetCode: retcode,
		Message: status,
		Data: fiber.Map{
			"message": message,
			"fonts":   fonts,
			"texts":   texts,
		},
	})
}
