package offices

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SelectUnitsParams struct {
	Operation int    `json:"operation"`
	Brcode    string `json:"brcode"`
	Staffid   string `json:"staffid"`
}

func GetUnits(c *fiber.Ctx) error {

	staffid := c.Params("id")
	GetUnitParameters := new(SelectUnitsParams)

	if err := c.BodyParser(&GetUnitParameters); err != nil {
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
	GetUnitParameters.Staffid = staffid

	result, err := Get_Units(GetUnitParameters)
	if err != nil {
		fmt.Println(err)
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}
	retCode := sharedfunctions.GetStringFromMap(result, "retCode")
	message := sharedfunctions.GetStringFromMap(result, "message")

	// Log operation
	logs.LOSLogs(c, GetRegionModule, staffid, retCode, message)
	return c.JSON(result)
}

type UpsertUnitsParams struct {
	Operation int    `json:"operation"`
	Brcode    string `json:"brcode"`
	UnitCode  int    `json:"unitCode"`
	UnitName  string `json:"unitName"`
	StaffID   string `json:"staffID"`
}

func UpsertUnits(c *fiber.Ctx) error {
	staffid := c.Params("id")
	upsertParameters := new(UpsertUnitsParams)
	params_select := new(SelectUnitsParams)

	if err := c.BodyParser(&upsertParameters); err != nil {
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

	result, err := Upsert_Units(staffid, upsertParameters, params_select)
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
	retCode := sharedfunctions.GetStringFromMap(result, "retCode")
	message := sharedfunctions.GetStringFromMap(result, "message")

	// Log operation
	logs.LOSLogs(c, GetRegionModule, staffid, retCode, message)
	return c.JSON(result)
}
