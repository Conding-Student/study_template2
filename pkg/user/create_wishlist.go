package users

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils/go-utils/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

type CreateWishlistRequest struct {
	StaffID     string
	WishList    string
	Description string
}

func CreateWishList(c *fiber.Ctx) error {
	db := database.DB
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

	if err := db.Raw("INSERT INTO public.wishlist (staff_id, wish, description) VALUES(?, ?, ?)", createWishList.StaffID, createWishList.WishList, createWishList.Description).Scan(&createWishList).Error; err != nil {
		log.Println(err)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to created wishlist.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Wishlist successfully submitted!",
		Data:    nil,
	})
}
