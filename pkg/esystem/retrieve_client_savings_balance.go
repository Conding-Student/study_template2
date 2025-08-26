package esystem

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type GetSavingsRequiredCreds struct {
	BrCode string `gorm:"not null"`
	Cid    int    `gorm:"not null"`
}

func GetClientSavingsBalance(c *fiber.Ctx) error {
	getSavingsCreds := new(GetSavingsRequiredCreds)

	if err := c.BodyParser(getSavingsCreds); err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	id := strconv.Itoa(getSavingsCreds.Cid)

	savingBalance, err := GetClientSavings(getSavingsCreds.Cid, getSavingsCreds.BrCode)
	if err != nil {
		logs.LOSLogs(c, LOSfeature, id, "500", err.Error()+" "+getSavingsCreds.BrCode)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "An error occured while fetching client savings to database.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	savingBalance.UpdateAsOf, err = sharedfunctions.FormatToDateOnly(savingBalance.UpdateAsOf)
	if err != nil {
		log.Println(err)
		logs.LOSLogs(c, LOSfeature, id, "500", err.Error()+" "+getSavingsCreds.BrCode)
	}

	logs.LOSLogs(c, LOSfeature, id, "200", "Client savings balance fetch successfully! "+getSavingsCreds.BrCode)
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Client savings balance fetch successfully!",
		Data:    savingBalance,
	})
}
