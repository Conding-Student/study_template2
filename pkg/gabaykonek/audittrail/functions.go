package audittrail

import (
	"chatbot/pkg/models/status"
	"chatbot/pkg/utils/go-utils/database"
)

func GetLogs(logsType string, operation int, startDate, endDate string) ([]map[string]any, int, string, string, string, error) {
	db := database.DB

	query := "SELECT * FROM logs.getlogs($1, $2, $3, $4)"

	var response []map[string]any
	if err := db.Raw(query, logsType, operation, startDate, endDate).Scan(&response).Error; err != nil {
		return nil, 500, "500", status.RetCode500, "Error fetching audit logs due a problem connecting to database!", err
	}

	return response, 200, "200", "Successful!", "Successfully fetch data.", nil
}
