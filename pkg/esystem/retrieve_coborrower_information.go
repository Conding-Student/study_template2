package esystem

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"

	"github.com/gofiber/fiber/v2"
)

func GetCoBorrowerInformationInEsystem(c *fiber.Ctx) error {
	staffid := c.Params("id")

	fetchedCoBorrowerInformation, err := GetCoBorrowerInformation(staffid)
	if err != nil {
		logs.LOSLogs(c, LOSfeature, staffid, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching co-borrower information list.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, LOSfeature, staffid, "200", ("Client Information retrieved successfully! "))
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Client Information retrieve successfully!",
		Data: fiber.Map{
			// "BrCode": brCode,
			"List": fetchedCoBorrowerInformation,
		},
	})
}
