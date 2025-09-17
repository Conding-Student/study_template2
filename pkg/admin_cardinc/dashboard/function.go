package dashboard

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
)

func FetchDashBoardData() (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT loan_application.get_dashboard_data()`).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_dashboard_data")

	return result, nil
}
