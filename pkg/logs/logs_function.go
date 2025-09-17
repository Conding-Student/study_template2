package logs

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
)

// insert logs from CA-Gabay
func LOSLogsInsert(userInput map[string]any) (map[string]any, error) {
	db := database.DB

	query := "SELECT * FROM logs.insert_loslog($1)"

	var response map[string]any
	if err := db.Raw(query, userInput).Scan(&response).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	return response, nil
}

// get logs from CA-Gabay
func GetLogsQuery(getLogs *GetCagabayLogs) (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw(`SELECT logs.get_loslogs(?)`, getLogs).Scan(&response).Error; err != nil {

		return nil, err

	}
	// Convert JSON string fields to proper JSON
	sharedfunctions.ConvertStringToJSONMap(response)

	result := sharedfunctions.GetMap(response, "get_loslogs")

	return result, nil
}

// insert audtrail
func InsertAuditTrail(input map[string]any) (map[string]any, error) {
	db := database.DB

	query := `SELECT logs.insert_audtrail($1)`

	var response map[string]any
	if err := db.Raw(query, input).Scan(&response).Error; err != nil {
		return nil, err
	}

	return response, nil
}

// Insert error log
func InsertErrorLog(input map[string]any) (map[string]any, error) {
	db := database.DB
	query := `SELECT logs.insert_errorlogs($1)`

	var response map[string]any
	if err := db.Raw(query, input).Scan(&response).Error; err != nil {
		return nil, err
	}

	return response, nil
}
