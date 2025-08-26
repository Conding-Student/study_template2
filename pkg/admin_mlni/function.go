package adminmlni

import (
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
)

func GetMlniStaffInfo() ([]map[string]any, error) {
	db := database.DB

	var cardStaffInfo []map[string]any
	if err := db.Raw("SELECT * FROM mlnitrackingfunc.getmlnistaff()").Scan(&cardStaffInfo).Error; err != nil {
		return nil, err
	}

	return cardStaffInfo, nil
}

func ManageMlniUser(request map[string]any) (map[string]any, bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM mlnitrackingfunc.updatemlnistaff($1)", request).Scan(&response).Error; err != nil {
		return nil, false, 500, "500", status.RetCode500, "An error occured while connecting to database.", err
	}

	sharedfunctions.ConvertStringToJSONMap(response)

	result := sharedfunctions.GetMap(response, "response")

	isSuccess := sharedfunctions.GetBoolFromMap(result, "issuccess")
	retcode := sharedfunctions.GetStringFromMap(result, "retcode")
	retCodeInt := sharedfunctions.GetIntFromMap(result, "retcode")
	status := sharedfunctions.GetStringFromMap(result, "status")
	message := sharedfunctions.GetStringFromMap(result, "message")

	if !isSuccess {
		return nil, isSuccess, retCodeInt, retcode, status, message, fmt.Errorf("there is a problem in updating user role")
	}

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nSuccessful: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	delete(result, "status")
	delete(result, "retcode")
	return result, isSuccess, retCodeInt, retcode, status, message, nil
}
