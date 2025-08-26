package kakfaconsumer

import "chatbot/pkg/utils/go-utils/database"

func StoreFailedDisbursement(kafkares map[string]any) (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.storedisbursementtransaction($1)", kafkares).Scan(&response).Error; err != nil {
		return nil, err
	}

	return response, nil
}
