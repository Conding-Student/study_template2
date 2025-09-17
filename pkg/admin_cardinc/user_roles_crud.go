package admincardinc

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

type AddUserRoles struct {
	Operation     int    `gorm:"not null"`
	Roles         string `gorm:"not null"`
	Min           int    `gorm:"not null"`
	Max           int    `gorm:"not null"`
	RecommendedBy string `gorm:"not null"`
	ApprovedBy    string `gorm:"not null"`
	Approver      bool   `gorm:"not null"`
	Description   string `gorm:"not null"`
}

func AddUserRole(c *fiber.Ctx) error {
	staffid := c.Params("id")

	userRoles := make(map[string]any)
	if err := c.BodyParser(&userRoles); err != nil {
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

	userRoles["staffid"] = staffid
	var resultData map[string]any
	resultData, err := Upsertgkroles(userRoles)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to fetch logs.",
				Error:     err,
			},
		})
	}

	return c.JSON(resultData)
}

func ViewUserRole(c *fiber.Ctx) error {
	staffid := c.Params("id")

	userRoles := make(map[string]any)
	if err := c.BodyParser(&userRoles); err != nil {
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

	userRoles["staffid"] = staffid
	var resultData map[string]any
	resultData, err := View_roles(userRoles)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to fetch logs.",
				Error:     err,
			},
		})
	}

	return c.JSON(resultData)
}

func DeleteuserRoles(c *fiber.Ctx) error {
	staffid := c.Params("id")

	userRoles := make(map[string]any)
	if err := c.BodyParser(&userRoles); err != nil {
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

	userRoles["staffid"] = staffid
	var resultData map[string]any
	resultData, err := Delete_gkrole(userRoles)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to fetch logs.",
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(resultData)
}

// func AddUserRole(c *fiber.Ctx) error {
// 	db := database.DB
// 	staffid := c.Params("id")

// 	userRoles := make(map[string]any)
// 	if err := c.BodyParser(&userRoles); err != nil {
// 		return c.Status(400).JSON(response.ResponseModel{
// 			RetCode: "400",
// 			Message: status.RetCode400,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to parse request",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	userRoles["staffid"] = staffid
// 	var resultData map[string]any
// 	if err := db.Raw("SELECT * FROM gabaykonekfunc.upsert_gkrole($1)", userRoles).Scan(&resultData).Error; err != nil {
// 		return c.Status(401).JSON(response.ResponseModel{
// 			RetCode: "401",
// 			Message: status.RetCode401,
// 			Data: errors.ErrorModel{
// 				Message:   "failed to create or update user roles",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	sharedfunctions.ConvertStringToJSONMap(resultData)
// 	result := sharedfunctions.GetMap(resultData, "result")

// 	return c.Status(200).JSON(result)
// }

// func ViewUserRole(c *fiber.Ctx) error {
// 	db := database.DB
// 	staffid := c.Params("id")

// 	userRoles := make(map[string]any)
// 	if err := c.BodyParser(&userRoles); err != nil {
// 		return c.Status(400).JSON(response.ResponseModel{
// 			RetCode: "400",
// 			Message: status.RetCode400,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to parse request",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	userRoles["staffid"] = staffid
// 	var resultData map[string]any
// 	if err := db.Raw("SELECT * FROM gabaykonekfunc.roles($1)", userRoles).Scan(&resultData).Error; err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Internal Server Error",
// 			Data: errors.ErrorModel{
// 				Message:   "Problem connecting to database",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	sharedfunctions.ConvertStringToJSONMap(resultData)
// 	result := sharedfunctions.GetMap(resultData, "result")

// 	return c.Status(200).JSON(result)
// }

// func DeleteuserRoles(c *fiber.Ctx) error {
// 	db := database.DB
// 	staffid := c.Params("id")

// 	userRoles := make(map[string]any)
// 	if err := c.BodyParser(&userRoles); err != nil {
// 		return c.Status(400).JSON(response.ResponseModel{
// 			RetCode: "400",
// 			Message: status.RetCode400,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to parse request",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	userRoles["staffid"] = staffid
// 	var resultData map[string]any
// 	if err := db.Raw("SELECT * FROM gabaykonekfunc.delete_gkrole($1)", userRoles).Scan(&resultData).Error; err != nil {
// 		return c.Status(401).JSON(response.ResponseModel{
// 			RetCode: "401",
// 			Message: status.RetCode401,
// 			Data: errors.ErrorModel{
// 				Message:   "failed to create or update user roles",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	sharedfunctions.ConvertStringToJSONMap(resultData)
// 	result := sharedfunctions.GetMap(resultData, "result")

// 	return c.Status(200).JSON(result)
// }
