package handler

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func InstitutionsAndEstablished(c *fiber.Ctx) error {
	db := database.DB

	var institutionsEstablished []InstitutionEstablished
	query := `SELECT * FROM public.institutions_established 
	ORDER BY 
		CAST(SUBSTRING(date_established FROM '[0-9]{4}') AS INTEGER) ASC,
		CASE 
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'January' THEN 1
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'February' THEN 2
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'March' THEN 3
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'April' THEN 4
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'May' THEN 5
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'June' THEN 6
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'July' THEN 7
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'August' THEN 8
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'September' THEN 9
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'October' THEN 10
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'November' THEN 11
			ELSE 12
		END ASC,
		CAST(SUBSTRING(date_established FROM '[0-9]{2}(?=\)$)') AS INTEGER) ASC;
	`
	if err := db.Raw(query).Scan(&institutionsEstablished).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch institutions and date established",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    institutionsEstablished,
	})
}

type Institution struct {
	InstitutionID int    `gorm:"column:id"`
	Institutions  string `gorm:"column:institutions"`
	InstiCode     string `gorm:"column:insti_code"`
	Description   string `gorm:"column:description"`
}

func Institutions(c *fiber.Ctx) error {
	db := database.DB

	var institutionsEstablished []Institution
	query := `SELECT id, institutions, insti_code, description  FROM public.institutions_established 
	ORDER BY 
		CAST(SUBSTRING(date_established FROM '[0-9]{4}') AS INTEGER) ASC,
		CASE 
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'January' THEN 1
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'February' THEN 2
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'March' THEN 3
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'April' THEN 4
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'May' THEN 5
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'June' THEN 6
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'July' THEN 7
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'August' THEN 8
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'September' THEN 9
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'October' THEN 10
			WHEN SUBSTRING(date_established FROM '\((\w+)') = 'November' THEN 11
			ELSE 12
		END ASC,
		CAST(SUBSTRING(date_established FROM '[0-9]{2}(?=\)$)') AS INTEGER) ASC;
	`
	if err := db.Raw(query).Scan(&institutionsEstablished).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch institutions list",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    institutionsEstablished,
	})
}
