package triviafacts

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RepoQueryParams struct {
	ID int64 `json:"id"`
}
type TriviaAndArticles struct {
	ID          int    `json:"id"`
	Featurename string `json:"featurename"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Image       string `json:"image"` // Remove form:"-" to accept JSON
	Link        string `json:"link"`
	Linktitle   string `json:"linktitle"`
	Author      string `json:"author"`
}

type DeletionArticleTrivia struct {
	Id int `json:"id"`
}

func GetArticles(c *fiber.Ctx) error {
	result, err := Get_Articles()
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch articles",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}

func GetTrivia(c *fiber.Ctx) error {
	result, err := Get_Trivia()
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

	return c.JSON(result)
}

func InsertArticleOrTrivia(c *fiber.Ctx) error {
	staffID := c.Params("id")
	article := new(TriviaAndArticles)

	// Parse JSON body
	if err := c.BodyParser(article); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid JSON data",
			Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
		})
	}

	// Handle base64 image
	if article.Image != "" {
		// Remove data URL prefix if present (e.g., "data:image/png;base64,")
		base64Data := article.Image
		if idx := strings.Index(base64Data, ","); idx != -1 {
			base64Data = base64Data[idx+1:]
		}

		fileBytes, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return c.Status(400).JSON(response.ResponseModel{
				RetCode: "400",
				Message: "Invalid base64 image",
				Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
			})
		}

		// Generate unique filename
		filename := fmt.Sprintf("image_%s_%d_%s.jpg", article.Featurename, time.Now().UnixNano(), staffID)

		// Upload to GitHub
		if githubURL, err := UploadToGitHub(filename, fileBytes); err == nil {
			article.Image = githubURL
		} else {
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Cannot upload image",
				Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
			})
		}
	}

	// Save to DB (pass as JSONB)
	result, err := Insert_ArticleOrTrivia(article)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to insert article",
			Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
		})
	}

	return c.JSON(result)
}

// DeleteArticleOrTrivia handles delete request
func DeleteArticleOrTrivia(c *fiber.Ctx) error {
	//staffID := c.Params("id")
	params := new(TriviaAndArticles)

	if err := c.BodyParser(&params); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400", Message: "Bad Request",
			Data: errors.ErrorModel{Message: "Failed to parse request", IsSuccess: false, Error: err},
		})
	}

	result, err := Get_FeatureImage(params)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500", Message: status.RetCode500,
			Data: errors.ErrorModel{Message: "Failed to insert article", IsSuccess: false, Error: err},
		})
	}

	if imageURL := sharedfunctions.GetStringFromMap(result, "featureImage"); imageURL != "" {
		if parts := strings.Split(imageURL, "/uploads/"); len(parts) == 2 {
			if err := DeleteFromGitHub(parts[1]); err != nil {
				return c.Status(500).JSON(response.ResponseModel{
					RetCode: "500", Message: "GitHub delete failed",
					Data: errors.ErrorModel{Message: "Cannot delete image from GitHub", IsSuccess: false, Error: err},
				})
			}
		}
	}

	resultDB, err := Delete_triviaorfacts(params)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500", Message: status.RetCode500,
			Data: errors.ErrorModel{Message: "Failed to insert article", IsSuccess: false, Error: err},
		})
	}
	return c.JSON(resultDB)
}

func UpdateArticleOrTrivia(c *fiber.Ctx) error {
	staffID := c.Params("id")
	trivia := new(TriviaAndArticles)

	// Parse JSON data
	if err := c.BodyParser(trivia); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid JSON data",
			Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
		})
	}

	// Handle base64 image if provided
	if trivia.Image != "" && !strings.HasPrefix(trivia.Image, "http") {
		// Remove data URL prefix if present (e.g., "data:image/png;base64,")
		base64Data := trivia.Image
		if idx := strings.Index(base64Data, ","); idx != -1 {
			base64Data = base64Data[idx+1:]
		}

		fileBytes, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return c.Status(400).JSON(response.ResponseModel{
				RetCode: "400",
				Message: "Invalid base64 image",
				Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
			})
		}

		// Generate filename
		filename := fmt.Sprintf("dashboard_image_%s_%d_%s.png", trivia.Featurename, time.Now().UnixNano(), staffID)

		// Upload/Update to GitHub
		if githubURL, err := UpdateFileOnGitHub(filename, fileBytes); err == nil {
			trivia.Image = githubURL
		} else {
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Cannot upload/update image",
				Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
			})
		}
	}

	// Get existing feature image
	imgLink, err := Get_FeatureImage(trivia)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data:    errors.ErrorModel{Message: "Failed to fetch existing feature image", IsSuccess: false, Error: err},
		})
	}

	// Delete old GitHub image if exists and a new image was uploaded
	if trivia.Image != "" {
		imageURL := sharedfunctions.GetStringFromMap(imgLink, "FeatureImage")
		if imageURL != "" {
			if parts := strings.Split(imageURL, "/uploads/"); len(parts) == 2 {
				if err := DeleteFromGitHub(parts[1]); err != nil {
					return c.Status(500).JSON(response.ResponseModel{
						RetCode: "500",
						Message: "GitHub delete failed",
						Data:    errors.ErrorModel{Message: "Cannot delete image from GitHub", IsSuccess: false, Error: err},
					})
				}
			}
		}
	}

	// Call DB update function
	result, err := Update_ArticlesOrTrivia(trivia)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to update trivia",
			Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
		})
	}

	return c.JSON(result)
}

// func UpdateArticles(c *fiber.Ctx) error {
// 	editArticles := new(EditTriviaAndArticles)

// 	if err := c.BodyParser(&editArticles); err != nil {
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

// 	// Call helper
// 	result, err := Update_Articles(editArticles)
// 	if err != nil {
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

// 	return c.JSON(result)
// }

// func GetArticles(c *fiber.Ctx) error {
// 	db := database.DB

// 	var trivia []map[string]any
// 	if err := db.Raw("SELECT * FROM GetArticles()").Scan(&trivia).Error; err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Internal server error",
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

// func UpdateArticles(c *fiber.Ctx) error {
// 	db := database.DB
// 	editArticles := new(EditTriviaAndArticles)

// 	if err := c.BodyParser(&editArticles); err != nil {
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
// 		SET title = ?, contents = ?, link_title = ?, link = ?, author = ?, feature_name = ?, feature_image = ?
// 		WHERE id = ?
// 	`
// 	if err := db.Exec(updateQuery, editArticles.Title, editArticles.Contents, editArticles.LinkTitle, editArticles.Link, editArticles.FeatureName, editArticles.FeatureImage, editArticles.ID).Error; err != nil {
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: status.RetCode500,
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to update article",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Success!",
// 		Data:    editArticles,
// 	})
// }
