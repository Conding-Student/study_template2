package users

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
)

func UpdateProfilePic(staffid, photo string) (bool, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM userprofile.updatephoto($1, $2)", staffid, photo).Scan(&response).Error; err != nil {
		return false, "An error occured while updating profile picture.", err
	}

	isSuccess := sharedfunctions.GetBoolFromMap(response, "issuccess")
	message := sharedfunctions.GetStringFromMap(response, "message")
	// errormesssage := sharedfunctions.GetStringFromMap(response, "errdetails")

	if !isSuccess {
		return isSuccess, message, fmt.Errorf("an error occured while updating profile picture")
	}

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nUpdate Photo successful: ", isSuccess)
	fmt.Println("Staff ID: ", staffid)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	return isSuccess, message, nil
}

func PinCreation(cagabayReqBody map[string]any) (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM userprofile.createpin($1)", cagabayReqBody).Scan(&response).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	response = sharedfunctions.GetMap(response, "createpin")

	return response, nil
}

func Insert_wishlist(params *CreateWishlistRequest) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT * FROM public.add_to_wishlist(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "add_to_wishlist")

	return result, nil
}
