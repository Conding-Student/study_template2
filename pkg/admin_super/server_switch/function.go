package switches

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
)

func Update_Switch(params *UpdateSwitchParams) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT public.update_switch(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "update_switch")

	return result, nil
}

func Get_Switch() (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Note: no params here since it's a fetch-all
	if err := db.Raw(`SELECT public.get_switch()`).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_switch")

	return result, nil
}
