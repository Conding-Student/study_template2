package sharedfunctions

import (
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
)

func SyncUserData(staffId string) (bool, int, string, string, string, error) {

	response, isCSuccess, cRetCodeInt, cRetCode, cResponseStatus, cResponseMessage, err := CreateUpdateUser(staffId)
	if err != nil {
		fmt.Println(err)
		return isCSuccess, cRetCodeInt, cRetCode, cResponseStatus, cResponseMessage, err
	}

	fmt.Println("Syncing response: ", response)

	return isCSuccess, cRetCodeInt, cRetCode, cResponseStatus, cResponseMessage, nil
}

func FetchProfile(accessType, staffID string) (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM userprofile.staffprofile($1, $2)", accessType, staffID).Scan(&response).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	ConvertStringToJSONMap(response)

	return response, nil
}

func CreateUpdateUser(staffId string) (map[string]any, bool, int, string, string, string, error) {
	db := database.DB

	var updateResult map[string]any
	staffInfo, err := Hcis(staffId)
	if err != nil {
		fmt.Println(err)
		// return nil, false, 500, "500", status.RetCode500, "An error occured while syncing your data.", err
		if err := db.Raw("SELECT * FROM userprofile.staffdatanotsync($1)", staffId).Scan(&updateResult).Error; err != nil {
			return nil, false, 500, "500", status.RetCode500, "An error occured while syncing your data.", err
		}
	} else {
		if err := db.Raw("SELECT * FROM userprofile.insertupdateprofile($1, $2)", true, staffInfo).Scan(&updateResult).Error; err != nil {
			return nil, false, 500, "500", status.RetCode500, "An error occured while syncing your data.", err
		}
	}

	ConvertStringToJSONMap(updateResult)

	result := GetMap(updateResult, "staffdata")
	// dataMap := GetMap(result, "data")
	// userData := GetMap(dataMap, "userData")
	// personalInfo := GetMap(userData, "personalInfo")

	// staffID := GetStringFromMap(personalInfo, "staffID")
	isSuccess := GetBoolFromMap(result, "issuccess")
	status := GetStringFromMap(result, "status")
	message := GetStringFromMap(result, "message")
	retCodeInt := GetIntFromMap(result, "retcodeint")
	retCode := GetStringFromMap(result, "retcode")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nStaffID/Username: ", staffId)
	fmt.Println("Successful Syncing/Creation: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return nil, isSuccess, retCodeInt, retCode, status, message, fmt.Errorf("an error occured while creating or syncing user data")
	}

	return updateResult, isSuccess, retCodeInt, retCode, status, message, nil

}

func FetchDesignation(staffId string) (string, error) {
	db := database.DB

	var designation string
	if err := db.Raw("SELECT * FROM userprofile.getdesignation($1)", staffId).Scan(&designation).Error; err != nil {
		return "Failed to identify user email.", err
	}

	return designation, nil
}

func FetchUserEmail(staffId string) (string, error) {
	db := database.DB

	var email string
	if err := db.Raw("SELECT * FROM userprofile.getemail(?)", staffId).Scan(&email).Error; err != nil {
		return "", err
	}

	return email, nil
}

func UpdateUser(request map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM public.update_user($1)AS update_user", request).Scan(&result).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Convert JSON string fields to proper JSON
	ConvertStringToJSONMap(result)

	results := GetMap(result, "update_user")

	return results, nil
}

func GetAllUsers() (map[string]any, error) {
	db := database.DB

	var cardStaffInfo map[string]any
	if err := db.Raw("SELECT * FROM userprofile.getalluser()").Scan(&cardStaffInfo).Error; err != nil {
		return nil, err
	}

	return cardStaffInfo, nil
}
