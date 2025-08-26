package offices

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

type SelectBranchesParams struct {
	Operation int
	Region    string
}

func GetBranches(c *fiber.Ctx) error {
	db := database.DB

	getBranchParameters := new(SelectBranchesParams)
	if err := c.BodyParser(&getBranchParameters); err != nil {
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

	var branches []map[string]any
	if err := db.Raw("SELECT * FROM cardincoffices.getbranches($1, $2)", getBranchParameters.Operation, getBranchParameters.Region).Scan(&branches).Error; err != nil {
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

	if len(branches) == 0 {
		branches = make([]map[string]any, 0)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successfully fetch branches",
		Data: fiber.Map{
			"offices": branches,
		},
	})
}

type UpsertBranchesParams struct {
	Operation  int
	Region     string
	Brcode     string
	BranchName string
	Active     bool
	StaffID    string
}

func UpsertBranches(c *fiber.Ctx) error {
	db := database.DB

	requestBody := make(map[string]any)
	if err := c.BodyParser(&requestBody); err != nil {
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

	var result map[string]any
	if err := db.Raw("SELECT * FROM cardincoffices.upsertbranches($1)", requestBody).Scan(&result).Error; err != nil {
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

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsertbranches")

	return c.JSON(result)
}
