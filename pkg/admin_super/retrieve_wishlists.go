package administrator

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils/go-utils/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetWishList(c *fiber.Ctx) error {
	db := database.DB

	var WishLists []map[string]any
	if err := db.Raw("SELECT * FROM getwishlist()").Scan(&WishLists).Error; err != nil {
		log.Println(err)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to fetch wishlist.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "WishList fetch successfully!",
		Data:    WishLists,
	})
}
