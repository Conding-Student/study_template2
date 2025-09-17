package offices

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

//	type UpsertClusterParams struct {
//		Operation int
//		Cluster   string
//		StaffID   string
//	}
type UpsertClusterParams struct {
	Operation int    `json:"operation"`
	Cluster   string `json:"cluster"`
	StaffID   string `json:"staffID"`
}

var GetClustersModule = "Clusters Module"

func GetClusters(c *fiber.Ctx) error {
	staffID := c.Params("id") // optional for logging

	// Delegate to query function
	result, err := Get_Clusters()
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
	logs.LOSLogs(c, GetClustersModule, staffID, retCode, message)

	return c.JSON(result)
}
func UpsertCluster(c *fiber.Ctx) error {
	staffID := c.Params("id") // optional for logging
	upsertParameters := new(UpsertClusterParams)
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

	result, err := Upsert_Cluster(upsertParameters)
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
	logs.LOSLogs(c, GetClustersModule, staffID, retCode, message)
	return c.JSON(result)
}

// func GetClusters(c *fiber.Ctx) error {
// 	db := database.DB

// 	var clusters []map[string]any
// 	if err := db.Raw("SELECT * FROM cardincoffices.getallclusters()").Scan(&clusters).Error; err != nil {
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

// 	if len(clusters) == 0 {
// 		clusters = make([]map[string]any, 0)
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Successfully fetch clusters",
// 		Data: fiber.Map{
// 			"offices": clusters,
// 		},
// 	})
// }

// func UpsertCluster(c *fiber.Ctx) error {
// 	db := database.DB

// 	upsertParameters := new(UpsertClusterParams)
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

// 	query := "SELECT * FROM cardincoffices.upsertcluster($1, $2, $3)"

// 	var existed bool
// 	if err := db.Raw(query, upsertParameters.Operation, upsertParameters.Cluster, upsertParameters.StaffID).Scan(&existed).Error; err != nil {
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
// 				Message:   "Cluster already exist. Please check clusters list.",
// 				IsSuccess: false,
// 				Error:     nil,
// 			},
// 		})
// 	}

// 	var result fiber.Map
// 	var message string
// 	switch upsertParameters.Operation {
// 	case 0:
// 		message = "Cluster successfully created."
// 		result = fiber.Map{
// 			"operation": "create",
// 			"request":   upsertParameters,
// 			"existed":   existed,
// 		}
// 	case 2:
// 		message = "Cluster successfully updated."
// 		result = fiber.Map{
// 			"operation": "update",
// 			"request":   upsertParameters,
// 		}
// 	case 3:
// 		message = "Cluster successfully deleted."
// 		result = fiber.Map{
// 			"operation": "delete",
// 			"request":   upsertParameters,
// 		}
// 	case 4:
// 		message = "Cluster reset successful."
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: message,
// 		Data:    result,
// 	})
// }
