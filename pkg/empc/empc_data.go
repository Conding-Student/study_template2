package empc

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type StaffInfoRequest struct {
	HcisId string `json:"hcisId"`
}

type StaffInfoResponse struct {
	FirstName     string `json:"firstName"`
	MiddleName    string `json:"middleName"`
	LastName      string `json:"lastName"`
	Birthday      string `json:"birthdate"`
	StaffId       string `json:"staffId"`
	Email         string `json:"email"`
	Mobile        string `json:"mobile"`
	EmpcCid       string `json:"empcCid"`
	PayrollStatus string `json:"payrollStatus"`
}

var divider = "-----------------------------------------------------------------------------------------------------------------"

func EmpcStaffInfo(c *fiber.Ctx) error {
	fmt.Println(divider)
	staffRequest := new(StaffInfoRequest)
	if err := c.BodyParser(staffRequest); err != nil {
		fmt.Println("retCode 401")
		fmt.Println("Invalid Request")
		fmt.Println("Failed to parse request", err.Error())
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

	staffId := staffRequest.HcisId

	responseBody, err := EmpcempcStaffInfo(staffId)
	if err != nil {
		logs.LOSLogs(c, EMPCFeature, staffId, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "We encountered an error while fetching your data in EMPC. Please try again later.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Check if the response contains an error field
	var empResponse map[string]any
	if err := json.Unmarshal(responseBody, &empResponse); err != nil {
		logs.LOSLogs(c, EMPCFeature, staffId, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "We encountered an error while fetching your data in EMPC. Please try again later.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	empcResponseData := new(StaffInfoResponse)
	actions := empResponse["content"].(map[string]any)["actions"].([]any)

	// Iterate over the actions array
	for _, action := range actions {
		actionMap := action.(map[string]any) // Convert action to a map[string]any

		// Extract the value for each field and assign it to the corresponding field in empcResponseData
		switch actionMap["field_name"].(string) {
		case "firstName":
			empcResponseData.FirstName = actionMap["value"].(string)
		case "middleName":
			empcResponseData.MiddleName = actionMap["value"].(string)
		case "lastName":
			empcResponseData.LastName = actionMap["value"].(string)
		case "birthDate":
			empcResponseData.Birthday = actionMap["value"].(string)
		case "staffID":
			empcResponseData.StaffId = actionMap["value"].(string)
		case "emailAddress":
			empcResponseData.Email = actionMap["value"].(string)
		case "mobile":
			empcResponseData.Mobile = actionMap["value"].(string)
		case "cid":
			empcResponseData.EmpcCid = actionMap["value"].(string)
		case "payrollStatus":
			empcResponseData.PayrollStatus = actionMap["value"].(string)
		}
	}

	if empcResponseData.EmpcCid == "0" {
		logs.LOSLogs(c, EMPCFeature, staffId, "404", fmt.Sprintf("The staff ID you provided %s was not found or is inactive. Please double-check your staff ID and try again.", staffId))
		return c.Status(404).JSON(response.ResponseModel{
			RetCode: "404",
			Message: status.RetCode404,
			Data: errors.ErrorModel{
				Message:   fmt.Sprintf("The staff ID you provided %s was not found or is inactive. Please double-check your staff ID and try again.", staffId),
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	empcBirthday, err := sharedfunctions.FormatToDateOnly(empcResponseData.Birthday)
	if err != nil {
		logs.LOSLogs(c, EMPCFeature, staffId, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   fmt.Sprintf("The staff ID you provided %s was not found or is inactive. Please double-check your staff ID and try again.", staffId),
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Handle the successful response
	empcResponse := StaffInfoResponse{
		FirstName:     empcResponseData.FirstName,
		MiddleName:    empcResponseData.MiddleName,
		LastName:      empcResponseData.LastName,
		Birthday:      empcBirthday,
		StaffId:       empcResponseData.StaffId,
		Email:         empcResponseData.Email,
		Mobile:        empcResponseData.Mobile,
		EmpcCid:       empcResponseData.EmpcCid,
		PayrollStatus: empcResponseData.PayrollStatus,
	}

	logs.LOSLogs(c, EMPCFeature, staffId, "200", "Successfully fetch data for EMPC. "+empcResponse.FirstName+" "+empcResponse.LastName+" "+empcResponse.PayrollStatus)
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    empcResponse,
	})
}
