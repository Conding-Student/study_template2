package audittrail

import (
	// "chatbot/pkg/models/errors"
	// "chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	// "strconv"
	// "github.com/gofiber/fiber/v2"
)

func GetLogs(logsType string, operation int, startDate, endDate string) ([]map[string]any, int, string, string, string, error) {
	db := database.DB

	query := "SELECT * FROM logs.getlogs($1, $2, $3, $4)"

	var response []map[string]any
	if err := db.Raw(query, logsType, operation, startDate, endDate).Scan(&response).Error; err != nil {
		return nil, 500, "500", status.RetCode500, "Error fetching audit logs due a problem connecting to database!", err
	}

	return response, 200, "200", "Successful!", "Successfully fetch data.", nil
}

func GetLogsTrial(userInput map[string]any) (map[string]any, error) {
	db := database.DB

	query := "SELECT * FROM logs.getlogstrial($1)"

	var response map[string]any
	if err := db.Raw(query, userInput).Scan(&response).Error; err != nil {
		return nil, err
	}

	// Convert any stringified JSON fields into nested maps/arrays
	sharedfunctions.ConvertStringToJSONMap(response)

	// Extract the actual function output (under "getlogstrial")
	result := sharedfunctions.GetMap(response, "getlogstrial")

	return result, nil
}

// // validateRequiredField checks if a field is provided and matches the expected type.
// // expectedType can be "string" or "int".
// func validateRequiredField(c *fiber.Ctx, fieldValue, fieldName, expectedType string) (any, error) {
// 	if fieldValue == "" {
// 		return nil, c.Status(400).JSON(response.ResponseModel{
// 			RetCode: "400",
// 			Message: "Missing required field: " + fieldName,
// 			Data: errors.ErrorModel{
// 				IsSuccess: false,
// 				Message:   "Invalid or missing " + fieldName,
// 				Error:     nil,
// 			},
// 		})
// 	}

// 	switch expectedType {
// 	case "string":
// 		// Return string directly
// 		return fieldValue, nil

// 	case "int":
// 		// Try to convert to int
// 		val, err := strconv.Atoi(fieldValue)
// 		if err != nil {
// 			return nil, c.Status(400).JSON(response.ResponseModel{
// 				RetCode: "400",
// 				Message: "Invalid " + fieldName,
// 				Data: errors.ErrorModel{
// 					IsSuccess: false,
// 					Message:   fieldName + " must be a valid integer",
// 					Error:     err,
// 				},
// 			})
// 		}
// 		return val, nil

// 	default:
// 		return nil, c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Unexpected type for field: " + fieldName,
// 			Data: errors.ErrorModel{
// 				IsSuccess: false,
// 				Message:   "Validator misconfigured for " + fieldName,
// 				Error:     nil,
// 			},
// 		})
// 	}
// }
