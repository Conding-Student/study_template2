package administrator

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func GetWishList(c *fiber.Ctx) error {

	//var WishLists map[string]any
	resultData, err := Retrieve_wishlist()
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to fetch logs.",
				Error:     err,
			},
		})
	}

	return c.JSON(resultData)
}
