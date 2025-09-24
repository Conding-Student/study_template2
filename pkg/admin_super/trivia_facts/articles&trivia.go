package triviafacts

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	//"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TriviaAndArticles struct {
	ID          int    `json:"id" form:"id"`
	Featurename string `json:"featurename" form:"featurename"`
	Title       string `json:"title" form:"title"`
	Content     string `json:"content" form:"content"`
	Image       string `json:"image" form:"-"`
	Link        string `json:"link" form:"link"`
	Linktitle   string `json:"linktitle" form:"linktitle"`
	Author      string `json:"author" form:"author"`
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

	// Parse form data
	if err := c.BodyParser(article); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid form data",
			Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
		})
	}

	// Handle image upload
	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(400).JSON(response.ResponseModel{
				RetCode: "400",
				Message: "Invalid image file",
				Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
			})
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Cannot read image file",
				Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
			})
		}

		// Validate required fields
		if article.Featurename == "" || article.Title == "" {
			return c.Status(400).JSON(response.ResponseModel{
				RetCode: "400",
				Message: "Featurename and Title are required",
				Data:    errors.ErrorModel{Message: "Featurename and Title are required", IsSuccess: false},
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

	// Save to DB
	result, err := Insert_ArticleOrTrivia(staffID, article)
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
	staffID := c.Params("id")
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

	resultDB, err := Delete_triviaorfacts(staffID, params)
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

	// Parse form data
	if err := c.BodyParser(trivia); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid form data",
			Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
		})
	}

	// Handle image upload if provided
	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(400).JSON(response.ResponseModel{
				RetCode: "400",
				Message: "Invalid image file",
				Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
			})
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Cannot read image file",
				Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
			})
		}

		// Validate required fields
		if trivia.Featurename == "" || trivia.Title == "" {
			return c.Status(400).JSON(response.ResponseModel{
				RetCode: "400",
				Message: "Bad Request",
				Data:    errors.ErrorModel{Message: "Featurename and Title are required", IsSuccess: false},
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

	imgLink, err := Get_FeatureImage(trivia)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data:    errors.ErrorModel{Message: "Failed to fetch existing feature image", IsSuccess: false, Error: err},
		})
	}

	// Delete old GitHub image if exists
	if imageURL := sharedfunctions.GetStringFromMap(imgLink, "FeatureImage"); imageURL != "" {
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

	// Call DB update function
	result, err := Update_ArticlesOrTrivia(staffID, trivia)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to update trivia",
			Data:    errors.ErrorModel{Message: err.Error(), IsSuccess: false, Error: err},
		})
	}

	return c.JSON(result)
}
