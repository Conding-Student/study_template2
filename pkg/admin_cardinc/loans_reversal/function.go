package loanreversal

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
)

func ReverseLoanApplication(userInput map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Call Postgres function
	if err := db.Raw("SELECT gabaykonekfunc.reverseloanapplication($1)", userInput).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert stringified JSON to map
	sharedfunctions.ConvertStringToJSONMap(result)

	// Correct key for GetMap
	result = sharedfunctions.GetMap(result, "reverseloanapplication")

	return result, nil
}
func loan_viewing(userInput map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Call Postgres function
	if err := db.Raw("SELECT gabaykonekfunc.getlistofdisloan($1)", userInput).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert stringified JSON to map
	sharedfunctions.ConvertStringToJSONMap(result)
	loans := sharedfunctions.GetMap(result, "getlistofdisloan")

	return loans, nil
}
