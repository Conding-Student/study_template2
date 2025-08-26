package logs

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GetCagabayLogs struct {
	StartDate  string
	EndDate    string
	StatusCode string
}

func LOSLogs(c *fiber.Ctx, feature, identification, statusCode, description string) error {
	db := database.DB

	insertQuery := `
	INSERT INTO logs.cagabay_logs
	(feature, identification, status_code, description, remarks, date)
	VALUES(?, ?, ?, ?, ?, ?)
	`

	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		log.Println(err)
		return err
	}

	manilaTime := time.Now().In(loc)
	logsTime := manilaTime.Format("2006-01-02 15:04:05.999999-07:00")

	remarks := ""

	if err := db.Exec(insertQuery, feature, identification, statusCode, description, remarks, logsTime).Error; err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Logs:", feature, identification, statusCode, description, remarks, logsTime)

	return nil
}

func GetLogs(c *fiber.Ctx) error {
	db := database.DB
	getLogs := new(GetCagabayLogs)

	if err := c.BodyParser(&getLogs); err != nil {
		log.Println(err)
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

	if getLogs.StatusCode == "" {
		getLogs.StatusCode = "200"
	}

	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		log.Println(err)
		return err
	}

	if getLogs.StartDate == "" {
		manilaTime := time.Now().In(loc)
		logsTime := manilaTime.Format("2006-01-02 15:04:05.999999-07:00")
		getLogs.StartDate = logsTime
	}

	if getLogs.EndDate == "" {
		parsedStartDate, err := time.ParseInLocation("2006-01-02 15:04:05.999999-07:00", getLogs.StartDate, loc)
		if err != nil {
			log.Println("Error parsing StartDate:", err)
			return c.Status(401).JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   "Failed to determine the date range.",
					IsSuccess: false,
					Error:     err,
				},
			})
		}
		endDate := time.Date(parsedStartDate.Year(), parsedStartDate.Month(), parsedStartDate.Day(), 0, 0, 0, 0, loc)
		getLogs.EndDate = endDate.Format("2006-01-02 15:04:05.999999-07:00")
	}

	// if getLogs.EndDate == "" {
	// 	endDate := time.Now().In(loc)
	// 	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, loc)
	// 	getLogs.EndDate = endDate.Format("2006-01-02 15:04:05.999999-07:00")
	// }

	fmt.Println("StartDate:", getLogs.StartDate)
	fmt.Println("EndDate:", getLogs.EndDate)

	getLogsQuery := `
		SELECT * FROM logs.cagabay_logs
		WHERE status_code = ? AND date > ? AND date < ?
		ORDER BY date DESC
	`

	var logs []map[string]any
	if err := db.Raw(getLogsQuery, getLogs.StatusCode, getLogs.StartDate, getLogs.EndDate).Scan(&logs).Error; err != nil {
		log.Println(err)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to fetch CA-GABAY Logs.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    logs,
	})
}
