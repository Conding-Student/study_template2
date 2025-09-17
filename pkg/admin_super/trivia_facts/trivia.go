package triviafacts

// import (
// 	"chatbot/pkg/models/errors"
// 	"chatbot/pkg/models/response"
// 	"chatbot/pkg/models/status"
// 	"chatbot/pkg/sharedfunctions"
// 	"encoding/base64"
// 	"fmt"
// 	"strings"

// 	"time"

// 	//"chatbot/pkg/utils/go-utils/database"

// 	"github.com/gofiber/fiber/v2"
// )

// func GetTrivia(c *fiber.Ctx) error {
// 	db := database.DB

// 	var trivia []map[string]any
// 	if err := db.Raw("SELECT * FROM GetTrivia()").Scan(&trivia).Error; err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: status.RetCode500,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to fetch secondary features",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Success!",
// 		Data:    trivia,
// 	})
// }

// func UpdateTrivia(c *fiber.Ctx) error {
// 	db := database.DB
// 	editTrivia := new(EditTriviaAndArticles)

// 	if err := c.BodyParser(&editTrivia); err != nil {
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

// 	updateQuery := `
// 		UPDATE public.secondary_features
// 		SET title = ?, contents = ?, link_title = ?, link = ?, feature_name = ?, feature_image = ?
// 		WHERE id = ?
// 	`
// 	if err := db.Exec(updateQuery, editTrivia.Title, editTrivia.Contents, editTrivia.LinkTitle, editTrivia.Link, editTrivia.FeatureName, editTrivia.FeatureImage, editTrivia.ID).Error; err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: status.RetCode500,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to update trivia",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Success!",
// 		Data:    editTrivia,
// 	})
// }
