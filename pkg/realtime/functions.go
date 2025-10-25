package realtime

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
)

func AddWhitelist(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw(`SELECT elli.addwhitelistws(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "addwhitelistws")

	return result, nil
}

func DelWhitelist(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw(`SELECT elli.deletewhitelistws(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "deletewhitelistws")

	return result, nil
}

func DelfeatureWhitelist(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw(`SELECT elli.deletefeaturewhitelistws(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "deletefeaturewhitelistws")

	return result, nil
}

func GetWhitelist(params map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT elli.getwhitelistws(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getwhitelistws")

	return result, nil
}

func GetfeatureWhitelist(params map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT elli.getfeaturewhitelistws(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getfeaturewhitelistws")

	return result, nil
}
