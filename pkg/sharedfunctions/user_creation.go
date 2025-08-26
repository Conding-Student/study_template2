package sharedfunctions

import (
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils"
	"chatbot/pkg/utils/go-utils/database"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"
	"fmt"
	"log"
)

func AccountCreationMobile(staffInfo map[string]any) (bool, int, string, string, string, error) {
	db := database.DB
	secretKey := utils.GetEnv("SECRET_KEY")

	firstName := GetStringFromMap(staffInfo, "firstName")
	middleName := GetStringFromMap(staffInfo, "middleName")
	lastName := GetStringFromMap(staffInfo, "lastName")
	temporaryUsername, err := CreateTemporaryUsername(firstName, middleName, lastName)
	if err != nil {
		return false, 401, "401", status.RetCode500, "Failed to generate temporary username. First and Last name is empty", err
	}

	// 2 generate password for new user
	genSuccess, genMessage, temporyPassword, err := GenerateSaveTempPass(2, "", "", "", "")
	if err != nil {
		return false, 401, "401", status.RetCode500, "Failed to generate temporary password", err
	}

	if !genSuccess {
		return genSuccess, 401, "401", status.RetCode500, genMessage, err
	}

	encryptedPassword, err := encryptDecrypt.Encrypt(temporyPassword, secretKey)
	if err != nil {
		return false, 401, "401", status.RetCode401, "Failed to encrypt temporary password", err
	}

	institution := GetStringFromMap(staffInfo, "institution")
	var institutionCode int
	if err := db.Raw("SELECT insti_code FROM public.institutions_established WHERE institutions = ?", institution).Scan(&institutionCode).Error; err != nil {
		return false, 500, "500", status.RetCode500, "Failed to fetch insti code", err
	}

	staffInfo["password"] = encryptedPassword
	staffInfo["username"] = temporaryUsername
	var regResult map[string]any
	if err := db.Raw("SELECT * FROM userprofile.usercreationmobile($1)", staffInfo).Scan(&regResult).Error; err != nil {
		log.Println(err)
		return false, 500, "500", status.RetCode500, "Failed to create user account due to a problem connecting to server", err
	}

	regResultSuccess := GetBoolFromMap(regResult, "issuccess")
	regResultRetCode := GetStringFromMap(regResult, "retcode")
	regResultRetCodeInt := GetIntFromMap(regResult, "retcode")
	regResultStatus := GetStringFromMap(regResult, "status")
	regResultMessage := GetStringFromMap(regResult, "message")

	emailAddress := GetStringFromMap(staffInfo, "email")
	fmt.Println(regResult)
	if regResultSuccess {
		if err := SendEmail(emailAddress, temporyPassword, temporaryUsername, 2); err != nil {
			return false, 401, "401", "Invalid Request!", "Failed to send account credentials.", err
		}
	} else {
		return regResultSuccess, regResultRetCodeInt, regResultRetCode, regResultStatus, regResultMessage, nil
	}

	// fmt.Println("Second Stage: ", staffInfo)

	responseBody, err := SaveToSoteriaGoMobile(staffInfo, temporaryUsername, temporyPassword, institutionCode)
	if err != nil {
		fmt.Println(err)
		// return false, 401, "401", "Invalid Request!", "Failed to send account credentials", err
	}

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nCreation successful: ", regResultSuccess)
	fmt.Println("Response from Soteria: ", responseBody)
	fmt.Println("Message: ", regResultMessage)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	return regResultSuccess, regResultRetCodeInt, regResultRetCode, regResultStatus, regResultMessage, nil
}

func AccountCreationAdmin(staffInfo map[string]any) (bool, int, string, string, string, error) {
	db := database.DB
	secretKey := utils.GetEnv("SECRET_KEY")

	firstName := GetStringFromMap(staffInfo, "firstName")
	middleName := GetStringFromMap(staffInfo, "middleName")
	lastName := GetStringFromMap(staffInfo, "lastName")
	temporaryUsername, err := CreateTemporaryUsername(firstName, middleName, lastName)
	if err != nil {
		return false, 401, "401", status.RetCode500, "Failed to generate temporary username. First and Last name is empty", err
	}

	// 2 generate password for new user
	genSuccess, genMessage, temporyPassword, err := GenerateSaveTempPass(2, "", "", "", "")
	if err != nil {
		return false, 401, "401", status.RetCode500, "Failed to generate temporary password", err
	}

	if !genSuccess {
		return genSuccess, 401, "401", status.RetCode500, genMessage, err
	}

	encryptedPassword, err := encryptDecrypt.Encrypt(temporyPassword, secretKey)
	if err != nil {
		return false, 401, "401", status.RetCode401, "Failed to encrypt temporary password", err
	}

	institution := GetStringFromMap(staffInfo, "institution")
	var institutionCode int
	if err := db.Raw("SELECT insti_code FROM public.institutions_established WHERE institutions = ?", institution).Scan(&institutionCode).Error; err != nil {
		return false, 500, "500", status.RetCode500, "Failed to fetch insti code", err
	}

	staffInfo["password"] = encryptedPassword
	staffInfo["username"] = temporaryUsername
	var regResult map[string]any
	if err := db.Raw("SELECT * FROM userprofile.usercreationadmin($1)", staffInfo).Scan(&regResult).Error; err != nil {
		log.Println(err)
		return false, 500, "500", status.RetCode500, "Failed to create user account due to a problem connecting to server", err
	}

	regResultSuccess := GetBoolFromMap(regResult, "issuccess")
	regResultRetCode := GetStringFromMap(regResult, "retcode")
	regResultRetCodeInt := GetIntFromMap(regResult, "retcode")
	regResultStatus := GetStringFromMap(regResult, "status")
	regResultMessage := GetStringFromMap(regResult, "message")

	emailAddress := GetStringFromMap(staffInfo, "email")

	// fmt.Println(regResult)
	if regResultSuccess {
		if err := SendEmail(emailAddress, temporyPassword, temporaryUsername, 2); err != nil {
			return false, 401, "401", "Invalid Request!", "Failed to send account credentials.", err
		}
	} else {
		return regResultSuccess, regResultRetCodeInt, regResultRetCode, regResultStatus, regResultMessage, nil
	}

	// fmt.Println("Second Stage: ", staffInfo)

	responseBody, err := SaveToSoteriaGoMobile(staffInfo, temporaryUsername, temporyPassword, institutionCode)
	if err != nil {
		fmt.Println(err)
		// return false, 401, "401", "Invalid Request!", "Failed to send account credentials", err
	}

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nCreation successful: ", regResultSuccess)
	fmt.Println("Response from Soteria: ", responseBody)
	fmt.Println("Message: ", regResultMessage)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	return regResultSuccess, regResultRetCodeInt, regResultRetCode, regResultStatus, regResultMessage, nil
}
