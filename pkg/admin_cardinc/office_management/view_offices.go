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

type AllBranch struct {
	BrCode     string `gorm:"column:brcode"`
	BranchName string `gorm:"column:branch_name"`
}

type UnitPerBranch struct {
	BrCode   string `gorm:"column:brcode"`
	UnitCode string `gorm:"column:unitcode"`
	UnitName string `gorm:"column:unit_name"`
}

type UnitPerBranchRequest struct {
	BrCode string `json:"brcode"`
}

func ViewCardIncBranches(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")

	params := make(map[string]any)
	// if err := c.BodyParser(&params); err != nil {
	// 	return c.Status(401).JSON(response.ResponseModel{
	// 		RetCode: "401",
	// 		Message: status.RetCode401,
	// 		Data: errors.ErrorModel{
	// 			Message:   "Failed to parse request",
	// 			IsSuccess: false,
	// 			Error:     err,
	// 		},
	// 	})
	// }

	params["cid"] = id
	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.officesgetbranches($1)", params).Scan(&result).Error; err != nil {
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
	result = sharedfunctions.GetMap(result, "response")

	return c.JSON(result)
}

func ViewCardIncUnits(c *fiber.Ctx) error {
	db := database.DB
	var brcode UnitPerBranchRequest
	if err := c.BodyParser(&brcode); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: status.RetCode400,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Query for approving authority records
	var unitList []UnitPerBranch
	if err := db.Raw("SELECT * FROM loan_application.unit($1)", &brcode.BrCode).Scan(&unitList).Error; err != nil {
		fmt.Println(err)
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

	// Return a success response with the approving authority data
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "List of units per branch retrieved successfully",
		Data:    unitList,
	})
}
