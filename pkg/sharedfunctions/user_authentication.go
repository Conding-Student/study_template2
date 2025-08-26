package sharedfunctions

import (
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
)

func SaveTokenToDB(staffid, token string) (bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM authentication.savetoken($1, $2)", staffid, token).Scan(&response).Error; err != nil {
		return false, 500, "500", status.RetCode500, "Login failed! Error storing token!", err
	}

	isSuccess := GetBoolFromMap(response, "issuccess")
	retCodeInt := GetIntFromMap(response, "retcodeint")
	retCode := GetStringFromMap(response, "retcode")
	status := GetStringFromMap(response, "status")
	message := GetStringFromMap(response, "message")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nStaffID/Username: ", staffid)
	fmt.Println("Storing token successful: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return isSuccess, retCodeInt, retCode, status, message, fmt.Errorf("login failed! Error storing token")
	}

	return isSuccess, retCodeInt, retCode, status, message, nil
}

func ValidateToken(token string) (bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM authentication.validatetoken($1)", token).Scan(&response).Error; err != nil {
		return false, 500, "500", status.RetCode500, "Failed to validate token!", err
	}

	isSuccess := GetBoolFromMap(response, "issuccess")
	retCodeInt := GetIntFromMap(response, "retcodeint")
	retCode := GetStringFromMap(response, "retcode")
	status := GetStringFromMap(response, "status")
	message := GetStringFromMap(response, "message")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nSesion is valid: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return isSuccess, retCodeInt, retCode, status, message, fmt.Errorf("failed to validate token")
	}

	return isSuccess, retCodeInt, retCode, status, message, nil
}
