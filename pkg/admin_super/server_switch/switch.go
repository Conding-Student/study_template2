package switches

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

type UpdateSwitchParams struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Switch  bool   `json:"switch"`
}

func GetSwitch(c *fiber.Ctx) error {
	result, err := Get_Switch()
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

func UpdateSwitch(c *fiber.Ctx) error {
	updateParameters := new(UpdateSwitchParams)
	if err := c.BodyParser(&updateParameters); err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	result, err := Update_Switch(updateParameters)
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

// type EditSwitch struct {
// 	ID      int
// 	Message string
// 	Switch  bool
// }

// func GetSwitch(c *fiber.Ctx) error {
// 	db := database.DB

// 	var serverSwitch []map[string]any
// 	if err := db.Raw("SELECT * FROM public.server_switch ORDER BY id ASC").Scan(&serverSwitch).Error; err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: status.RetCode500,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to fetch server switches.",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Successful",
// 		Data:    serverSwitch,
// 	})
// }

// func UpdateSwitch(c *fiber.Ctx) error {
// 	db := database.DB
// 	editSwitch := new(EditSwitch)

// 	if err := c.BodyParser(&editSwitch); err != nil {
// 		return c.Status(401).JSON(response.ResponseModel{
// 			RetCode: "401",
// 			Message: status.RetCode401,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to parse request",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	// Validate required fields
// 	if editSwitch.ID == 0 || editSwitch.Message == "" {
// 		return c.Status(401).JSON(response.ResponseModel{
// 			RetCode: "401",
// 			Message: status.RetCode401,
// 			Data: errors.ErrorModel{
// 				Message:   "ID and Message are required fields",
// 				IsSuccess: false,
// 				Error:     nil,
// 			},
// 		})
// 	}

// 	updateQuery := `
// 		UPDATE public.server_switch
// 		SET message = ?, switch = ?
// 		WHERE id = ?
// 	`

// 	if err := db.Exec(updateQuery, editSwitch.Message, editSwitch.Switch, editSwitch.ID).Error; err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: status.RetCode500,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to fetch server switches.",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Server switch updated successfully",
// 		Data:    editSwitch,
// 	})
// }
