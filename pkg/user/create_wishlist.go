package users

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

type CreateWishlistRequest struct {
	StaffID     string `json:"staff_id"`
	Wish        string `json:"wish"`
	Description string `json:"description"`
}

func CreateWishList(c *fiber.Ctx) error {

	createWishList := new(CreateWishlistRequest)

	if err := c.BodyParser(&createWishList); err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}
	result, err := Insert_wishlist(createWishList)
	if err != nil {
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

	return c.JSON(result)
}
