package administrator

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
)

func Retrieve_wishlist() (map[string]any, error) {
	db := database.DB
	var result map[string]any

	err := db.Raw("SELECT public.get_wishlist()").Scan(&result).Error
	if err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_wishlist")
	return result, nil
}
