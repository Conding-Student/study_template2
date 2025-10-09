package admincardinc

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"

	"github.com/gofiber/fiber/v2"
)

func GetCardIncUsers(c *fiber.Ctx) error {
	staffInfo, err := GetCardIncStaffInfo()
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(staffInfo)
}

// func GetCardIncUsers(c *fiber.Ctx) error {

// 	cardStaffInfo, err := GetCardIncStaffInfo()
// 	if err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: status.RetCode500,
// 			Data: errors.ErrorModel{
// 				Message:   "Problem connecting to database",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Successful!",
// 		Data:    cardStaffInfo,
// 	})
// }
