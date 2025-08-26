package loans

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

type jsonBRequestBody map[string]any

func InsertLoanApplication(c *fiber.Ctx) error {
	loanCreation := make(jsonBRequestBody)

	if err := c.BodyParser(&loanCreation); err != nil {
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

	loanCreation["source"] = c.Get("sourceApplication")

	qrDetails, isSuccess, retCodeInt, retCode, lStatus, message, err := LoanCreation(loanCreation)
	if err != nil {
		fmt.Println(err)
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: lStatus,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	refCode := sharedfunctions.GetStringFromMap(qrDetails, "refcode")
	qrCode, err := qrcode.Encode(refCode, qrcode.High, 256)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to generate QR code",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	loanInfo := fiber.Map{
		"referenceCode": refCode,
		"qrCode":        qrCode,
		"timeCreated":   sharedfunctions.GetStringFromMap(qrDetails, "date"),
	}

	return c.Status(retCodeInt).JSON(response.ResponseModel{
		RetCode: retCode,
		Message: "Loan application created successfully",
		Data:    loanInfo,
	})
}
