package offices

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

type SelectCentersParams struct {
	Operation int
	Brcode    string
	UnitCode  int
}

func GetCenters(c *fiber.Ctx) error {
	db := database.DB

	getCenterParameters := new(SelectCentersParams)
	if err := c.BodyParser(&getCenterParameters); err != nil {
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

	var centers []map[string]any
	if err := db.Raw("SELECT * FROM cardincoffices.getcenters($1, $2, $3)", getCenterParameters.Operation, getCenterParameters.Brcode, getCenterParameters.UnitCode).Scan(&centers).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem conecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	if len(centers) == 0 {
		centers = make([]map[string]any, 0)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successfully fetch centers",
		Data: fiber.Map{
			"offices": centers,
		},
	})
}

type UpsertCentersParams struct {
	Operation  int
	Brcode     string
	UnitCode   int
	CenterCode string
	CenterName string
	StaffID    string
}

func UpsertCenters(c *fiber.Ctx) error {
	db := database.DB

	upsertParameters := new(UpsertCentersParams)
	if err := c.BodyParser(&upsertParameters); err != nil {
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

	query := "SELECT * FROM cardincoffices.upsertcenters($1, $2, $3, $4, $5, $6)"

	var existed bool
	if err := db.Raw(
		query,
		upsertParameters.Operation,
		upsertParameters.Brcode,
		upsertParameters.UnitCode,
		upsertParameters.CenterCode,
		upsertParameters.CenterName,
		upsertParameters.StaffID,
	).Scan(&existed).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem conecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	if existed {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Center code already exist. Please check center list.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	var result fiber.Map
	var message string
	switch upsertParameters.Operation {
	case 0:
		message = "Center successfully created."
		result = fiber.Map{
			"operation": "create",
			"request":   upsertParameters,
			"existed":   existed,
		}
	case 2:
		message = "Center successfully updated."
		result = fiber.Map{
			"operation": "update",
			"request":   upsertParameters,
		}
	case 3:
		message = "Center successfully deleted."
		result = fiber.Map{
			"operation": "delete",
			"request":   upsertParameters,
		}
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: message,
		Data:    result,
	})
}
