package admincardinc

import "chatbot/pkg/utils/go-utils/database"

func GetCardIncStaffInfo() ([]map[string]any, error) {
	db := database.DB

	var cardStaffInfo []map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getcardincstaff()").Scan(&cardStaffInfo).Error; err != nil {
		return nil, err
	}

	return cardStaffInfo, nil
}
