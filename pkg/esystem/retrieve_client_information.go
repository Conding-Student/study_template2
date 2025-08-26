package esystem

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetClientInformationInEsystem(c *fiber.Ctx) error {
	clientCreds := make(map[string]any)

	if err := c.BodyParser(&clientCreds); err != nil {
		log.Println(err)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request body.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	cid := c.Params("id")
	cidStr := sharedfunctions.GetStringFromMap(clientCreds, "Cid")
	brCode := sharedfunctions.GetStringFromMap(clientCreds, "Brcode")

	if cid != cidStr {
		logs.LOSLogs(c, LOSfeature, cidStr, "401", "Cid is invalid. Please double-check and try again. "+brCode)
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Cid is invalid. Please double-check and try again.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	fetchedClientInformation, retCode, err := GetClientInformationV2(clientCreds)
	if err != nil {
		logs.LOSLogs(c, LOSfeature, cidStr, retCode, err.Error()+" "+brCode)
		return c.JSON(fetchedClientInformation)
	}

	logs.LOSLogs(c, LOSfeature, cidStr, retCode, "Client Information retrieved successfully! "+brCode)
	return c.JSON(fetchedClientInformation)
}
