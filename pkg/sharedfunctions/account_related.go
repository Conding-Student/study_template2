package sharedfunctions

import (
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
	"net/url"
	"strconv"

	jwtToken "chatbot/pkg/utils/go-utils/fiber"

	"github.com/gofiber/fiber/v2"
)

func UpdatePassword(staffID, currentPassword, newPassword string) (bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM userprofile.updatepassword($1, $2, $3)", staffID, currentPassword, newPassword).Scan(&response).Error; err != nil {
		fmt.Println(err)
		return true, 500, "500", status.RetCode500, "An error occurred while updating your password.", err
	}

	isSuccess := GetBoolFromMap(response, "issuccess")
	retCode := GetStringFromMap(response, "retcode")
	retCodeInt, err := strconv.Atoi(retCode)
	if err != nil {
		fmt.Println(err)
		return true, 401, "401", status.RetCode401, "An error occurred while updating your password.", err
	}
	status := GetStringFromMap(response, "status")
	responseMessage := GetStringFromMap(response, "message")

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nChange Password successful: ", isSuccess)
	fmt.Println("Message: ", responseMessage)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return isSuccess, retCodeInt, retCode, status, responseMessage, fmt.Errorf("an error occurred while updating your password")
	}

	return true, retCodeInt, retCode, status, responseMessage, nil
}

func UpdateDevice(deviceInfo map[string]any) (bool, int, string, string, string, error) {
	db := database.DB
	db = db.Debug()

	var response map[string]any
	if err := db.Raw("SELECT * FROM userprofile.updatedevice($1)", deviceInfo).Scan(&response).Error; err != nil {
		fmt.Println(err)
		return true, 500, "500", status.RetCode500, "An error occurred while updating your device.", err
	}

	fmt.Println(GetStringFromMap(deviceInfo, "staffID"))
	fmt.Println(GetStringFromMap(deviceInfo, "deviceID"))

	isSuccess := GetBoolFromMap(response, "issuccess")
	retCode := GetStringFromMap(response, "retcode")
	retCodeInt, err := strconv.Atoi(retCode)
	if err != nil {
		fmt.Println(err)
		return true, 401, "401", status.RetCode401, "An error occurred while updating your device.", err
	}
	status := GetStringFromMap(response, "status")
	responseMessage := GetStringFromMap(response, "message")

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nChange Device successful: ", isSuccess)
	fmt.Println("Message: ", responseMessage)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return isSuccess, retCodeInt, retCode, status, responseMessage, fmt.Errorf("an error occurred while updating your device")
	}

	return true, retCodeInt, retCode, status, responseMessage, nil
}

func AccountLoginAdmin(staffid string, loginCreds map[string]any) (map[string]any, bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM public.accountloginvalidationadmin($1)", loginCreds).Scan(&response).Error; err != nil {
		return nil, false, 500, "500", status.RetCode500, "An error occured while validating your credentials.", err
	}

	isSuccess := GetBoolFromMap(response, "issuccess")
	retCode := GetStringFromMap(response, "retcode")
	retCodeInt := GetIntFromMap(response, "retcode")
	responseStatus := GetStringFromMap(response, "status")
	responseMessage := GetStringFromMap(response, "message")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nStaffID/Username: ", staffid)
	fmt.Println("Valid credentials: ", isSuccess)
	fmt.Println("Message: ", responseMessage)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return nil, isSuccess, retCodeInt, retCode, responseStatus, responseMessage, fmt.Errorf("an error occured while validating your credentials")
	}

	staffid = GetStringFromMap(response, "staffid")
	updateResult, isCSuccess, cRetCodeInt, cRetCode, cResponseStatus, cResponseMessage, err := CreateUpdateUser(staffid)
	if err != nil {
		fmt.Println(err)
		return nil, isCSuccess, cRetCodeInt, cRetCode, cResponseStatus, cResponseMessage, err
	}

	data := GetMap(updateResult, "staffdata")
	userData := GetMap(data, "data")
	// fmt.Println("User Data: ", userData)

	employmentInfo := GetMap(userData, "employmentInfo")

	accessToken, err := jwtToken.GenerateJWTSignedString(fiber.Map{
		"staffID":     staffid,
		"rolename":    GetStringFromMap(userData, "role"),
		"institution": GetStringFromMap(employmentInfo, "institution"),
		"area":        GetStringFromMap(employmentInfo, "area"),
		"unit":        GetStringFromMap(employmentInfo, "unit"),
		"designation": GetStringFromMap(employmentInfo, "designation"),
	})

	if err != nil {
		fmt.Println("An error occured while generating token.", err)
		return nil, false, 500, "500", status.RetCode500, "An error occured while generating token.", err
	}

	userData["accessToken"] = accessToken
	delete(userData, "enabledFeatures")

	isTSuccess, retTCodeInt, retTCode, tstatus, tmessage, err := SaveTokenToDB(staffid, accessToken)
	if err != nil {
		return nil, isTSuccess, retTCodeInt, retTCode, tstatus, tmessage, err
	}

	return userData, true, 200, "200", "Login successful!", "User successfully login.", nil
}

func AccountLoginV2(loginCreds map[string]any, staffid, deviceid string) (map[string]any, bool, int, string, string, string, error) {
	db := database.DB

	decodedDeviceID, err := url.QueryUnescape(deviceid)
	if err != nil {
		fmt.Println("Error identifying device id: ", err)
		return nil, false, 401, "401", status.RetCode401, "Failed to identify device id.", err
	}

	loginCreds["deviceid"] = decodedDeviceID
	var response map[string]any
	if err := db.Raw("SELECT * FROM public.accountloginvalidationv2($1)", loginCreds).Scan(&response).Error; err != nil {
		fmt.Println("An error occured while validating credentials: ", err)
		return nil, false, 500, "500", status.RetCode500, "An error occured while validating your credentials.", err
	}

	isSuccess := GetBoolFromMap(response, "issuccess")
	retCode := GetStringFromMap(response, "retcode")
	retCodeInt := GetIntFromMap(response, "retcode")
	responseStatus := GetStringFromMap(response, "status")
	responseMessage := GetStringFromMap(response, "message")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nStaffID/Username: ", staffid)
	fmt.Println("Valid credentials: ", isSuccess)
	fmt.Println("Message: ", responseMessage)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	if !isSuccess {
		fmt.Println("Error message: ", responseMessage)
		return nil, isSuccess, retCodeInt, retCode, responseStatus, responseMessage, fmt.Errorf("invalid credentials or account is locked")
	}

	staffid = GetStringFromMap(response, "staffid")
	updateResult, isCSuccess, cRetCodeInt, cRetCode, cResponseStatus, cResponseMessage, err := CreateUpdateUser(staffid)
	if err != nil {
		fmt.Println(cResponseMessage, ": ", err)
		return nil, isCSuccess, cRetCodeInt, cRetCode, cResponseStatus, cResponseMessage, err
	}

	result := GetMap(updateResult, "staffdata")
	data := GetMap(result, "data")
	userData := GetMap(data, "userData")
	// fmt.Println("User Data: ", userData)

	employmentInfo := GetMap(userData, "employmentInfo")

	// rolename := GetStringFromMap(data, "role")
	// institution := GetStringFromMap(employmentInfo, "institution")
	// instiCode := GetIntFromMap(employmentInfo, "instiCode")
	// area := GetStringFromMap(employmentInfo, "area")
	// unit := GetStringFromMap(employmentInfo, "unit")
	// designation := GetStringFromMap(employmentInfo, "designation")

	// fmt.Println(instiCode)
	// fmt.Println(rolename)
	// fmt.Println(institution)
	// fmt.Println(area)
	// fmt.Println(unit)
	// fmt.Println(designation)

	accessToken, err := jwtToken.GenerateJWTSignedString(fiber.Map{
		"staffID":     staffid,
		"rolename":    GetStringFromMap(data, "role"),
		"institution": GetStringFromMap(employmentInfo, "institution"),
		"insticode":   GetIntFromMap(employmentInfo, "instiCode"),
		"area":        GetStringFromMap(employmentInfo, "area"),
		"unit":        GetStringFromMap(employmentInfo, "unit"),
		"designation": GetStringFromMap(employmentInfo, "designation"),
	})

	if err != nil {
		fmt.Println("An error occured while generating token.", err)
		return nil, false, 500, "500", status.RetCode500, "An error occured while generating token.", err
	}

	data["accessToken"] = accessToken

	isTSuccess, retTCodeInt, retTCode, tstatus, tmessage, err := SaveTokenToDB(staffid, accessToken)
	if err != nil {
		return nil, isTSuccess, retTCodeInt, retTCode, tstatus, tmessage, err
	}

	return data, true, 200, "200", "Login successful!", "User successfully login.", nil
}

func AccountLogout(params map[string]any) (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM public.accountlogoutv2($1)", params).Scan(&response).Error; err != nil {
		fmt.Println("An error occured while logging out: ", err)
		return response, err
	}

	ConvertStringToJSONMap(response)
	data := GetMap(response, "response")

	return data, nil
}

func PinValidation(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM userprofile.validatepin($1)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	ConvertStringToJSONMap(result)
	result = GetMap(result, "validatepin")

	return result, nil
}
