package audittrail

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

type LogsRequest struct {
	Operation int
	StartDate string
	EndDate   string
}

func AccessLogs(c *fiber.Ctx) error {
	requestLog := new(LogsRequest)
	logsType := c.Get("logsType")

	var logs string
	if logsType == "0" {
		logs = "Access Logs"
	} else if logsType == "1" {
		logs = "System Logs"
	} else {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Invalid logs type.",
				Error:     nil,
			},
		})
	}

	if err := c.BodyParser(&requestLog); err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "Failed to parse request.",
				Error:     err,
			},
		})
	}

	auditLogs, retCodeInt, retCode, status, message, err := GetLogs(logsType, requestLog.Operation, requestLog.StartDate, requestLog.EndDate)
	if err != nil {
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: status,
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   message,
				Error:     err,
			},
		})
	}

	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: status,
		Data: fiber.Map{
			"message":  message,
			"logsType": logs,
			"result":   auditLogs,
		},
	})
}
