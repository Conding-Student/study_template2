package users

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/model"
	"chatbot/pkg/models/response"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// UpdateUser function handles the update operation for a user
func UpdateUser(c *fiber.Ctx) error {
	db := database.DB
	var requestBody map[string]any
	if err := c.BodyParser(&requestBody); err != nil {
		fmt.Println("retCode 301")
		fmt.Println("Invalid Data")
		fmt.Println("Failed to parsed data")
		return c.Status(301).JSON(response.ResponseModel{
			RetCode: "301",
			Message: "Invalid Data",
			Data: errors.ErrorModel{
				Message:   "Invalid input.\nPlease review your data.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	userID, ok := requestBody["id"].(string)
	if !ok || userID == "" {
		fmt.Println("retCode 404")
		fmt.Println("Bad Request")
		fmt.Println("Invalid user ID")
		return c.Status(404).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "Invalid data in the request body",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	password, ok := requestBody["password"].(string)
	if !ok || password == "" {
		fmt.Println("retCode 404")
		fmt.Println("Bad Request")
		fmt.Println("Invalid Password")
		return c.Status(404).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "Invalid data in the request body",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	// Find the user in the database
	var user model.UserSignUp
	result := db.First(&user, "id = ?", userID)
	if result.Error != nil {
		fmt.Println("retCode 404")
		fmt.Println("Bad Request")
		fmt.Println("The user did not found")
		return c.Status(404).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "User not found",
				IsSuccess: false,
				Error:     result.Error,
			},
		})
	}

	// Compare currentPassword with the existing password in the database
	if password != user.Password {
		fmt.Println("retCode 404")
		fmt.Println("Bad Request")
		fmt.Println("The user submit an invalid password")
		// Passwords don't match
		return c.Status(404).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "It seems that you have provided an invalid password. Kindly enter your correct password to successfully update your information.",
				IsSuccess: false,
				Error:     result.Error,
			},
		})
	}

	// Parse and bind the request body to the user struct
	if err := c.BodyParser(&user); err != nil {
		fmt.Println("retCode 404")
		fmt.Println("Bad Request")
		fmt.Println("Failed to parse the data in model.UserSignUp")
		return c.Status(404).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Bad request",
			Data: errors.ErrorModel{
				Message:   "Invalid input.\nPlease review your data.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Update the user in the database
	if err := db.Save(&user).Error; err != nil {
		fmt.Println("retCode 500")
		fmt.Println("Internal Server Error")
		fmt.Println("Could not update user info")
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Could not update user",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Your profile is successfully updated.",
		Data:    user,
	})
}
