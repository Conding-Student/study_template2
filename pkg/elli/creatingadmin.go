package elli

import (
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils"
	"chatbot/pkg/utils/go-utils/database"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"
	"fmt"
	"log"
)

func creating_admin_account(staffInfo map[string]any) (bool, int, string, string, string, error) {
	db := database.DB
	secretKey := utils.GetEnv("SECRET_KEY")

	// Sanitize and get user input
	firstName := sharedfunctions.GetStringFromMap(staffInfo, "firstName")
	middleName := sharedfunctions.GetStringFromMap(staffInfo, "middleName")
	lastName := sharedfunctions.GetStringFromMap(staffInfo, "lastName")

	// Create temporary username
	temporaryUsername, err := sharedfunctions.CreateTemporaryUsername(firstName, middleName, lastName)
	if err != nil {
		return false, 401, "401", status.RetCode500, "Failed to generate temporary username. First and Last name is empty", err
	}

	// Generate password
	genSuccess, genMessage, temporyPassword, err := sharedfunctions.GenerateSaveTempPass(2, "", "", "", "")
	if err != nil {
		return false, 401, "401", status.RetCode500, "Failed to generate temporary password", err
	}
	if !genSuccess {
		return genSuccess, 401, "401", status.RetCode500, genMessage, err
	}

	// Encrypt password
	encryptedPassword, err := encryptDecrypt.Encrypt(temporyPassword, secretKey)
	if err != nil {
		return false, 401, "401", status.RetCode401, "Failed to encrypt temporary password", err
	}

	// Place in map
	staffInfo["password"] = encryptedPassword
	staffInfo["username"] = temporaryUsername

	var regResult map[string]any
	if err := db.Raw("SELECT * FROM elli.usercreation($1)", staffInfo).Scan(&regResult).Error; err != nil {
		log.Println(err)
		return false, 500, "500", status.RetCode500, "Failed to create user account due to a problem connecting to server", err
	}

	regResultSuccess := sharedfunctions.GetBoolFromMap(regResult, "issuccess")
	regResultRetCode := sharedfunctions.GetStringFromMap(regResult, "retcode")
	regResultRetCodeInt := sharedfunctions.GetIntFromMap(regResult, "retcode")
	regResultStatus := sharedfunctions.GetStringFromMap(regResult, "status")
	regResultMessage := sharedfunctions.GetStringFromMap(regResult, "message")

	fmt.Println("Second Stage: ", staffInfo)
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nCreation successful: ", regResultSuccess)
	fmt.Println("Response from Soteria: ", regResult)
	fmt.Println("Message: ", regResultMessage)
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

	return regResultSuccess, regResultRetCodeInt, regResultRetCode, regResultStatus, regResultMessage, nil
}
