package features

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
)

func GetSecFeatures() (map[string]any, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM public.getsecfeat()").Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getsecfeat")

	return result, nil
}
