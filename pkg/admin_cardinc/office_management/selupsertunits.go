package offices

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetUnits(c *fiber.Ctx) error {
	db := database.DB
	staffid := c.Params("id")

	params := make(map[string]any)
	if err := c.BodyParser(&params); err != nil {
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

	params["staffid"] = staffid
	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.officesgetunits($1)", params).Scan(&result).Error; err != nil {
		fmt.Println(err)
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "officesgetunits")

	return c.JSON(result)
}

type UpsertUnitsParams struct {
	Operation int
	Brcode    string
	UnitCode  int
	UnitName  string
	StaffID   string
}

func UpsertUnits(c *fiber.Ctx) error {
	db := database.DB

	upsertParameters := new(UpsertUnitsParams)
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

	query := "SELECT * FROM cardincoffices.upsertunits($1, $2, $3, $4, $5)"

	var existed bool
	if err := db.Raw(query, upsertParameters.Operation, upsertParameters.Brcode, upsertParameters.UnitCode, upsertParameters.UnitName, upsertParameters.StaffID).Scan(&existed).Error; err != nil {
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
				Message:   "Unit code already exist. Please check unit list.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	var result fiber.Map
	var message string
	switch upsertParameters.Operation {
	case 0:
		message = "Unit successfully created."
		result = fiber.Map{
			"operation": "create",
			"request":   upsertParameters,
			"existed":   existed,
		}
	case 2:
		message = "Unit successfully updated."
		result = fiber.Map{
			"operation": "update",
			"request":   upsertParameters,
		}
	case 3:
		message = "Unit successfully deleted."
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
