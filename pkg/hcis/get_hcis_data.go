package hcis

import (
	"chatbot/pkg/authentication"
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type StaffInfoRequest struct {
	Operation int    `gorm:"not null"`
	Otp       string `gorm:"not null"`
	StaffId   string `gorm:"not null"`
}

type HcisInfoResponse struct {
	Cluster            string `json:"cluster,omitempty"`
	LastName           string `json:"lastName,omitempty"`
	ZipCode            string `json:"zipCode,omitempty"`
	FirstName          string `json:"firstName,omitempty"`
	ProvAddress        string `json:"provAddress,omitempty"`
	DateHired          string `json:"dateHired,omitempty"`
	Gender             string `json:"gender,omitempty"`
	EmploymentStatus   string `json:"employmentStatus,omitempty"`
	DateRegularization string `json:"dateRegularization,omitempty"`
	CivilStatus        string `json:"civilStatus,omitempty"`
	Institution        string `json:"institution,omitempty"`
	EmailAddress       string `json:"email,omitempty"`
	BrgyAddress        string `json:"brgyAddress,omitempty"`
	CityAddress        string `json:"cityAddress,omitempty"`
	Department         string `json:"department,omitempty"`
	StaffID            string `json:"staffId,omitempty"`
	Area               string `json:"area,omitempty"`
	EmploymentType     string `json:"employmentType,omitempty"`
	LocationAssignment string `json:"locationAssignment,omitempty"`
	NickName           string `json:"nickName,omitempty"`
	BirthDate          string `json:"birthdate,omitempty"`
	JobLevel           string `json:"jobLevel,omitempty"`
	JobPosition        string `json:"jobPosition,omitempty"`
	Unit               string `json:"unit,omitempty"`
	JobGrade           string `json:"jobGrade,omitempty"`
	MobilePhone        string `json:"mobile,omitempty"`
	MiddleName         string `json:"middleName,omitempty"`
	Region             string `json:"region,omitempty"`
	Age                string `json:"age,omitempty"`
	Status             string `json:"status,omitempty"`
}

var accountCreationFeature = "AccountCreation"

func HcisStaffInfo(c *fiber.Ctx) error {
	staffRequest := new(StaffInfoRequest)
	staffID := c.Params("id")

	if err := c.BodyParser(staffRequest); err != nil {
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

	staffId := staffRequest.StaffId
	if staffID != staffId {
		logs.LOSLogs(c, accountCreationFeature, staffId, "401", "It seems that your Staff ID is invalid. Please double check and try again.")
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "It seems that your Staff ID is invalid. Please double check and try again.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	if staffRequest.Operation == 2 {
		message, verified := authentication.VerifyOTP(staffRequest.Otp, staffRequest.StaffId, 0)
		if !verified {
			return c.Status(401).JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   message,
					IsSuccess: false,
					Error:     nil,
				},
			})
		}
	}

	hcisUserInfo, err := sharedfunctions.Hcis(staffId)
	if err != nil {
		logs.LOSLogs(c, accountCreationFeature, staffId, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "We encountered an error while fetching your data. Please try again later.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	staffInfo, ok := hcisUserInfo["StaffInfo"].([]interface{})
	if !ok || len(staffInfo) == 0 {
		logs.LOSLogs(c, accountCreationFeature, staffId, "404", fmt.Sprintf("The staff ID %s you've provided was not found. Please double-check and try again.", staffId))
		return c.Status(404).JSON(response.ResponseModel{
			RetCode: "404",
			Message: status.RetCode404,
			Data: errors.ErrorModel{
				Message:   fmt.Sprintf("The staff ID %s you've provided was not found. Please double-check and try again.", staffId),
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	// hcisInfo := staffInfo[0].(map[string]any)
	var hcisInfo map[string]any
	if len(staffInfo) > 0 {
		if info, ok := staffInfo[0].(map[string]any); ok {
			hcisInfo = info
		} else {
			logs.LOSLogs(c, accountCreationFeature, staffId, "404", fmt.Sprintf("The staff ID %s you've provided was not found. Please double-check and try again.", staffId))
			return c.Status(404).JSON(response.ResponseModel{
				RetCode: "404",
				Message: status.RetCode404,
				Data: errors.ErrorModel{
					Message:   fmt.Sprintf("The staff ID %s you've provided was not found. Please double-check and try again.", staffId),
					IsSuccess: false,
					Error:     nil,
				},
			})
		}
	} else {
		logs.LOSLogs(c, accountCreationFeature, staffId, "404", fmt.Sprintf("The staff ID %s you've provided was not found. Please double-check and try again.", staffId))
		return c.Status(404).JSON(response.ResponseModel{
			RetCode: "404",
			Message: status.RetCode404,
			Data: errors.ErrorModel{
				Message:   fmt.Sprintf("The staff ID %s you've provided was not found. Please double-check and try again.", staffId),
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	staffData := new(HcisInfoResponse)

	switch staffRequest.Operation {
	case 0:

		// staffData.StaffID = hcisInfo["staffID"].(string)
		// staffData.BirthDate = hcisInfo["birthDate"].(string)
		// staffData.MobilePhone = hcisInfo["mobilePhone"].(string)
		// staffStatus := hcisInfo["status"].(string)
		staffData.StaffID = sharedfunctions.GetStringFromMap(hcisInfo, "staffID")
		staffData.BirthDate = sharedfunctions.GetStringFromMap(hcisInfo, "birthDate")
		staffData.MobilePhone = sharedfunctions.GetStringFromMap(hcisInfo, "mobilePhone")
		staffStatus := sharedfunctions.GetStringFromMap(hcisInfo, "status")

		if staffStatus != "active" {
			logs.LOSLogs(c, accountCreationFeature, staffId, "404", fmt.Sprintf("It looks like the staff ID %s you provided has a status of %s. Please double-check and try again.", staffId, staffStatus))
			return c.Status(404).JSON(response.ResponseModel{
				RetCode: "404",
				Message: status.RetCode404,
				Data: errors.ErrorModel{
					Message:   fmt.Sprintf("It looks like the staff ID %s you provided has a status of %s. Please double-check and try again.", staffId, staffStatus),
					IsSuccess: false,
					Error:     nil,
				},
			})
		}

	case 1:

		// email := hcisInfo["emailAddress"].(string)
		email := sharedfunctions.GetStringFromMap(hcisInfo, "emailAddress")
		if message, err := authentication.GenerateOtp(email, staffId, 0); err != nil {
			logs.LOSLogs(c, accountCreationFeature, staffId, "401", "We encounter an error while sending a One Time Password (OTP) in your email. Please try again later. "+staffId)
			return c.Status(401).JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   message,
					IsSuccess: false,
					Error:     nil,
				},
			})
		}

	case 2, 3:

		staffData.StaffID = sharedfunctions.GetStringFromMap(hcisInfo, "staffID")
		staffData.BirthDate = sharedfunctions.GetStringFromMap(hcisInfo, "birthDate")
		staffData.MobilePhone = sharedfunctions.GetStringFromMap(hcisInfo, "mobilePhone")
		staffData.Status = sharedfunctions.GetStringFromMap(hcisInfo, "status")
		staffData.FirstName = sharedfunctions.GetStringFromMap(hcisInfo, "firstname")
		staffData.MiddleName = sharedfunctions.GetStringFromMap(hcisInfo, "middleName")
		staffData.LastName = sharedfunctions.GetStringFromMap(hcisInfo, "lastName")
		staffData.EmailAddress = sharedfunctions.GetStringFromMap(hcisInfo, "emailAddress")
		staffData.Cluster = sharedfunctions.GetStringFromMap(hcisInfo, "cluster")
		staffData.ZipCode = sharedfunctions.GetStringFromMap(hcisInfo, "zipCode")
		staffData.ProvAddress = sharedfunctions.GetStringFromMap(hcisInfo, "provAddress")
		staffData.DateHired = sharedfunctions.GetStringFromMap(hcisInfo, "dateHired")
		staffData.Gender = sharedfunctions.GetStringFromMap(hcisInfo, "gender")
		staffData.EmploymentStatus = sharedfunctions.GetStringFromMap(hcisInfo, "employmentStatus")
		staffData.DateRegularization = sharedfunctions.GetStringFromMap(hcisInfo, "dateRegularization")
		staffData.CivilStatus = sharedfunctions.GetStringFromMap(hcisInfo, "civilStatus")
		staffData.Institution = sharedfunctions.GetStringFromMap(hcisInfo, "institution")
		staffData.BrgyAddress = sharedfunctions.GetStringFromMap(hcisInfo, "brgyAddress")
		staffData.CityAddress = sharedfunctions.GetStringFromMap(hcisInfo, "cityAddress")
		staffData.Department = sharedfunctions.GetStringFromMap(hcisInfo, "department")
		staffData.Area = sharedfunctions.GetStringFromMap(hcisInfo, "area")
		staffData.EmploymentType = sharedfunctions.GetStringFromMap(hcisInfo, "employmentType")
		staffData.LocationAssignment = sharedfunctions.GetStringFromMap(hcisInfo, "locationAssignment")
		staffData.NickName = sharedfunctions.GetStringFromMap(hcisInfo, "nickName")
		staffData.JobLevel = sharedfunctions.GetStringFromMap(hcisInfo, "jobLevel")
		staffData.JobPosition = sharedfunctions.GetStringFromMap(hcisInfo, "jobPosition")
		staffData.Unit = sharedfunctions.GetStringFromMap(hcisInfo, "unit")
		staffData.JobGrade = sharedfunctions.GetStringFromMap(hcisInfo, "jobGrade")
		staffData.Region = sharedfunctions.GetStringFromMap(hcisInfo, "region")
		staffData.Age = sharedfunctions.GetStringFromMap(hcisInfo, "age")

	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    staffData,
	})
}
