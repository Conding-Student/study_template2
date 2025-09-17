// BY Norman Villegas
package offices

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func GetStaffName(c *fiber.Ctx) error {

	// Parse request JSON
	var reqBody map[string]any
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request!",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request",
				Error:     err,
			},
		})
	}

	var result map[string]any
	result, err := Get_fullname(reqBody)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "failed to full name",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}
func GetStaffByDesignation(c *fiber.Ctx) error {
	var reqBody map[string]any
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request!",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request",
				Error:     err,
			},
		})
	}

	result, err := GetStaffByDesignationDB(reqBody)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Problem connecting to database",
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}

func GetCenterByStaffID(c *fiber.Ctx) error {
	var reqBody map[string]any
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request!",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request",
				Error:     err,
			},
		})
	}

	result, err := GetCenterByStaffIDDB(reqBody)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Problem connecting to database",
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}

func UpdateCenterTagStaff(c *fiber.Ctx) error {
	var reqBody map[string]any
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request!",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request",
				Error:     err,
			},
		})
	}

	result, err := UpdateCenterStaffDB(reqBody)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Problem connecting to database",
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(result)
}

// func GetStaffName(c *fiber.Ctx) error {
// 	db := database.DB

// 	var request struct {
// 		StaffId string `json:"staffId"`
// 	}

// 	if err := c.BodyParser(&request); err != nil {
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

// 	var result map[string]any
// 	if err := db.Raw("SELECT * FROM userprofile.getfullname($1)", request.StaffId).Scan(&result).Error; err != nil {
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

// 	sharedfunctions.ConvertStringToJSONMap(result)

// 	responseData := sharedfunctions.GetMap(result, "response")
// 	isSuccess := sharedfunctions.GetBoolFromMap(responseData, "issuccess")
// 	status := sharedfunctions.GetStringFromMap(responseData, "status")
// 	retCode := sharedfunctions.GetStringFromMap(responseData, "retcode")
// 	retCodeInt := sharedfunctions.GetIntFromMap(responseData, "retcode")
// 	message := sharedfunctions.GetStringFromMap(responseData, "message")
// 	staffName := sharedfunctions.GetMap(responseData, "data")

// 	if !isSuccess {
// 		return c.Status(retCodeInt).JSON(response.ResponseModel{
// 			RetCode: retCode,
// 			Message: status,
// 			Data: errors.ErrorModel{
// 				IsSuccess: isSuccess,
// 				Message:   message,
// 				Error:     nil,
// 			},
// 		})
// 	}

// 	return c.Status(retCodeInt).JSON(response.ResponseModel{
// 		RetCode: retCode,
// 		Message: status,
// 		Data:    staffName,
// 	})
// }

// func GetStaffByDesignation(c *fiber.Ctx) error {
// 	db := database.DB
// 	var staffList []map[string]any
// 	var request struct {
// 		OPERATION int `json:"operation"`
// 	}

// 	if err := c.BodyParser(&request); err != nil {
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

// 	// Execute the query
// 	if err := db.Raw("SELECT * FROM cardincoffices.getstaffbydesignation(?)", request.OPERATION).Scan(&staffList).Error; err != nil {
// 		// Optional: Log the error for internal tracking
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

// 	// Return successful response with the staffList data
// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Successful!",
// 		Data:    staffList,
// 	})
// }

// func GetCenterByStaffID(c *fiber.Ctx) error {
// 	db := database.DB

// 	var centerList []map[string]any

// 	var request struct {
// 		STAFFID string `json:"staff_id"`
// 	}

// 	if err := c.BodyParser(&request); err != nil {
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

// 	if err := db.Raw("SELECT * FROM cardincoffices.get_center_by_staff_id(?)", request.STAFFID).Scan(&centerList).Error; err != nil {
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

// 	var successMessage string

// 	if len(centerList) == 0 {
// 		successMessage = "No Center Tag For this Staff!"
// 	} else {
// 		successMessage = "Successful!"
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: successMessage,
// 		Data:    centerList,
// 	})
// }

// func UpdateCenterTagStaff(c *fiber.Ctx) error {
// 	db := database.DB

// 	var request struct {
// 		StaffId    string `json:"staff_id"`
// 		BrCode     string `json:"brcode"`
// 		UnitCode   int    `json:"unitcode"`
// 		CenterCode string `json:"centercode"`
// 	}

// 	if err := c.BodyParser(&request); err != nil {
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

// 	query := `SELECT cardincoffices.update_center_staffid($1, $2, $3, $4)`

// 	if err := db.Exec(query, request.StaffId, request.BrCode, request.UnitCode, request.CenterCode).Error; err != nil {
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
// 		Data:    nil,
// 	})
// }
