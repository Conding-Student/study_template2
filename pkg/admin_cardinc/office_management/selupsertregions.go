package offices

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

type SelectRegionsParams struct {
	SelectOption int    `json:"operation"`
	Cluster      string `json:"cluster"`
}

var GetRegionModule = "Region Module"

func GetRegions(c *fiber.Ctx) error {
	staffID := c.Params("id") // optional for logging
	getRegionParameters := new(SelectRegionsParams)

	// ✅ Parse body into struct
	if err := c.BodyParser(getRegionParameters); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}
	// Delegate to query function
	result, err := Get_Region(getRegionParameters)
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

	retCode := sharedfunctions.GetStringFromMap(result, "retCode")
	message := sharedfunctions.GetStringFromMap(result, "message")

	// Log operation
	logs.LOSLogs(c, GetRegionModule, staffID, retCode, message)

	return c.JSON(result)
}

type UpsertRegionParams struct {
	Operation int    `json:"operation"`
	Cluster   string `json:"cluster"`
	Region    string `json:"region"`
	StaffID   string `json:"staffID"`
}

func UpsertRegion(c *fiber.Ctx) error {
	staffid := c.Params("id")
	upsertParameters := new(UpsertRegionParams)
	selectparameters := new(SelectRegionsParams)
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

	result, err := Upsert_Region(staffid, upsertParameters, selectparameters)
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
	retCode := sharedfunctions.GetStringFromMap(result, "retCode")
	message := sharedfunctions.GetStringFromMap(result, "message")

	// Log operation
	logs.LOSLogs(c, GetRegionModule, staffid, retCode, message)
	return c.JSON(result)
}

// func GetRegions(c *fiber.Ctx) error {
// 	db := database.DB

// 	getRegionParameters := new(SelectRegionsParams)
// 	if err := c.BodyParser(&getRegionParameters); err != nil {
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

// 	var regions []map[string]any
// 	if err := db.Raw("SELECT * FROM cardincoffices.getregions($1, $2)", getRegionParameters.Operation, getRegionParameters.Cluster).Scan(&regions).Error; err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: status.RetCode500,
// 			Data: errors.ErrorModel{
// 				Message:   "Problem conecting to database",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	if len(regions) == 0 {
// 		regions = make([]map[string]any, 0)
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Successfully fetch regions",
// 		Data: fiber.Map{
// 			"offices": regions,
// 		},
// 	})
// }

// type UpsertRegionParams struct {
// 	Operation int
// 	Cluster   string
// 	Region    string
// 	StaffID   string
// }

// func UpsertRegion(c *fiber.Ctx) error {
// 	db := database.DB

// 	upsertParameters := new(UpsertRegionParams)
// 	if err := c.BodyParser(&upsertParameters); err != nil {
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

// 	query := "SELECT * FROM cardincoffices.upsertregion($1, $2, $3, $4)"

// 	var existed bool
// 	if err := db.Raw(query, upsertParameters.Operation, upsertParameters.Cluster, upsertParameters.Region, upsertParameters.StaffID).Scan(&existed).Error; err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: status.RetCode500,
// 			Data: errors.ErrorModel{
// 				Message:   "Problem conecting to database",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	if existed {
// 		return c.Status(401).JSON(response.ResponseModel{
// 			RetCode: "401",
// 			Message: status.RetCode401,
// 			Data: errors.ErrorModel{
// 				Message:   "Region already exist. Please check regions list.",
// 				IsSuccess: false,
// 				Error:     nil,
// 			},
// 		})
// 	}

// 	var result fiber.Map
// 	var message string
// 	switch upsertParameters.Operation {
// 	case 0:
// 		message = "Region successfully created."
// 		result = fiber.Map{
// 			"operation": "create",
// 			"request":   upsertParameters,
// 			"existed":   existed,
// 		}
// 	case 2:
// 		message = "Region successfully updated."
// 		result = fiber.Map{
// 			"operation": "update",
// 			"request":   upsertParameters,
// 		}
// 	case 3:
// 		message = "Region successfully deleted."
// 		result = fiber.Map{
// 			"operation": "delete",
// 			"request":   upsertParameters,
// 		}
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: message,
// 		Data:    result,
// 	})
// }
