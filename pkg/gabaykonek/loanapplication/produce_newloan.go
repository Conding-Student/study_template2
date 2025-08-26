package loans

import (
	"bytes"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

func produceToKafka(loanApplication map[string]any) (map[string]any, bool, int, string, string, string, error) {
	db := database.DB
	loanHeader := sharedfunctions.GetMap(loanApplication, "loanApp")
	clientInfo := sharedfunctions.GetMap(loanApplication, "customerInfo")
	// loanCoumputation := sharedfunctions.GetMap(loanApplication, "computation")
	// loanProductDetails := sharedfunctions.GetMap(loanApplication, "loanProducts")
	// demographic := sharedfunctions.GetMap(loanApplication, "demographic")
	// loanInventory := sharedfunctions.GetMap(loanApplication, "loanInventory")
	// coBorrowers := sharedfunctions.GetMap(loanApplication, "coborrowers")
	// midas := sharedfunctions.GetMap(loanApplication, "midas")
	// compliance := sharedfunctions.GetMap(loanApplication, "compliance")

	// staffID := sharedfunctions.GetStringFromMap(loanApplication, "staffid")
	// fmt.Println(staffID)

	// approvedAmount := sharedfunctions.GetFloatFromMap(loanHeader, "amount_approved")
	// frequency := sharedfunctions.GetIntFromMap(loanHeader, "payment_mode")
	// loanTerm := sharedfunctions.GetIntFromMap(loanHeader, "term")

	kafkaBaseUrl, err := sharedfunctions.GetBaseUrl(6)
	if err != nil {
		return nil, false, 500, "500", status.RetCode500, "An error occured while connecting to external API.", err
	}

	// loanProductCode := sharedfunctions.GetIntFromMap(loanProductDetails, "loan_product_code")

	// lrf, err := GetLoanLrf(approvedAmount, frequency, loanTerm)
	// if err != nil {
	// 	return nil, false, 500, "500", status.RetCode500, "An error occured while computing loan lrf.", err
	// }

	var approvedAmount float64
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getprincipalamount($1)", loanHeader).Scan(&approvedAmount).Error; err != nil {
		return nil, false, 500, "500", status.RetCode500, "An error occured while identifying loan details (cash or non-cash).", err
	}

	loanHeader["amount_approved"] = approvedAmount

	// fmt.Println("Loan Header: ", loanHeader)
	// fmt.Println("Loan Computation: ", loanCoumputation)
	// fmt.Println("Client Info: ", clientInfo)
	// fmt.Println("Loan Product Details: ", loanProductDetails)
	// fmt.Println("Demographic Survey: ", demographic)
	// fmt.Println("Loan Inventory: ", loanInventory)
	// fmt.Println("Co-Borrowers: ", coBorrowers)
	// fmt.Println("Midas: ", midas)
	// fmt.Println("Compliance: ", compliance)

	env := os.Getenv("ENVIRONMENT")
	var url string
	switch env {
	case "PROD":
		brcode := sharedfunctions.GetStringFromMap(clientInfo, "branch_code")
		url = kafkaBaseUrl + "/topic/newloan-" + brcode
	case "UAT":
		url = kafkaBaseUrl + "/topic/test.newloan-uat"
	default:
		return nil, false, 500, "500", status.RetCode500, "ENVIRONMENT is not set", err
	}

	newLoan := fiber.Map{
		"key":   "",
		"value": loanApplication,
	}

	requestBody, err := json.Marshal(newLoan)
	if err != nil {
		return nil, false, 500, "500", status.RetCode500, "An error occured while preparing loan details.", err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, false, 500, "500", status.RetCode500, "An error occured while forwarding loan details to external API.", err
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, false, 500, "500", status.RetCode500, "An error occured while forwarding loan details to external API.", err
	}

	var responseModel map[string]any
	if err := json.Unmarshal(body, &responseModel); err != nil {
		fmt.Println(err)
		return nil, false, 500, "500", status.RetCode500, "An error occured while identifying response from external API.", err
	}

	var success bool
	if response.StatusCode == 200 {
		success = sharedfunctions.GetBoolFromMap(responseModel, "success")
	} else {
		success = false
	}

	if !success {
		fmt.Println("--------------------------------------------------------------------------------------------------")
		fmt.Println("\nLoan is successfully submitted to kafka: ", success)
		fmt.Println("Status Code ", response.StatusCode)
		fmt.Println("Response Form kafka: ", responseModel)
		fmt.Println("\n------------------------------------------------------------------------------------------------")
		return nil, success, 500, "500", status.RetCode500, "Failed to forward loan details to external API.", fmt.Errorf("failed to forward loan details to external API")
	}

	return newLoan, true, 200, "200", "Success!", "Loan has been sent to kafka", nil
}
