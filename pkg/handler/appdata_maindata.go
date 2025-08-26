package handler

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func API(c *fiber.Ctx) error {
	db := database.DB
	var appVersion int
	userAgent := c.Get("User-Agent")

	appVersionStr := c.Params("appVersion")
	if appVersionStr == "" {
		appVersion = 0
	} else {
		var err error
		appVersion, err = strconv.Atoi(appVersionStr)
		if err != nil {
			logs.LOSLogs(c, "AppOpenning", appVersionStr, "401", err.Error()+" Device used: "+userAgent)
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   status.RetCode401,
					IsSuccess: false,
					Error:     err,
				},
			})
		}
	}

	params := make(map[string]any)
	params["appVersion"] = appVersion
	params["userAgent"] = userAgent
	var result map[string]any
	if err := db.Raw("SELECT * FROM public.maindata($1)", params).Scan(&result).Error; err != nil {
		logs.LOSLogs(c, "AppOpenning", appVersionStr, "500", err.Error()+" Device used: "+userAgent)
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "We couldn't identify the version of your CA-GABAY application. Please check for available updates on the Play Store and update your application.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "maindata")

	logs.LOSLogs(c, "AppOpenning", appVersionStr, "200", "Successfully fetch api. Device used: "+userAgent)
	return c.JSON(result)
}

func APIEndPoints(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("deviceid")

	var initialData map[string]any
	if err := db.Raw("SELECT * FROM public.getappdata()").Scan(&initialData).Error; err != nil {
		logs.LOSLogs(c, "AppOpenning", id, "500", err.Error())
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to fetchapp data.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(initialData)
	initialData = sharedfunctions.GetMap(initialData, "getappdata")

	logs.LOSLogs(c, "AppOpenning", id, "200", "Successfully fetch initial data.")
	return c.JSON(initialData)
}
