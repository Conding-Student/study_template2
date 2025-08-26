package sharedfunctions

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"

	"github.com/gofiber/fiber/v2"
)

func SendErrorResponse(
	c *fiber.Ctx,
	statusCode int, // HTTP status code (e.g. 400, 401, 500)
	retCode string, // Custom return code (string)
	message string, // Top-level messagea
	dataMessage string, // Data.Message (detailed)
	isSuccess bool, // Data.IsSuccess (usually false for errors)
	err error, // Actual error
) error {
	return c.Status(statusCode).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: message,
		Data: errors.ErrorModel{
			Message:   dataMessage,
			IsSuccess: isSuccess,
			Error:     err,
		},
	})
}

// SendErrorResponse (already exists)

// SendSuccessResponse sends a standardized success response
func SendSuccessResponse(
	c *fiber.Ctx,
	statusCode int, // HTTP status code (200, 201, etc.)
	retCode string, // Custom return code
	message string, // Success message
	data interface{}, // Any payload you want to return
) error {
	return c.Status(statusCode).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: message,
		Data:    data,
	})
}
