package creditline

import (
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
)

var module = "Credit Line"

type Float64 struct {
	Value float64
}

func (f Float64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", f.Value)), nil
}

func GetCredLineFields() (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.crdtlnegetallfields()").Scan(&response).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	response = sharedfunctions.GetMap(response, "crdtlnegetallfields")

	return response, nil
}

func GetCredLineProperties() (map[string]string, map[string]Float64, int, string, string, string, error) {
	db := database.DB

	var texts []map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.creditlinetexts()").Scan(&texts).Error; err != nil {
		return nil, nil, 500, "500", status.RetCode500, "Error fetching credit line fields due to problem connecting to database.", err
	}

	textMap := make(map[string]string)

	for _, text := range texts {
		idStr := fmt.Sprintf("%d", text["id"]) // Convert `id` to string
		textMap[idStr] = sharedfunctions.GetStringFromMap(text, "creditlinetext")
	}

	var fonts []map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.creditlinesizes()").Scan(&fonts).Error; err != nil {
		return nil, nil, 500, "500", status.RetCode500, "Error fetching credit line fields due to problem connecting to database.", err
	}

	fontsMap := make(map[string]Float64)

	for _, font := range fonts {
		fontsMap[sharedfunctions.GetStringFromMap(font, "sections")] = Float64{sharedfunctions.GetFloatFromMap(font, "font")}
	}

	return textMap, fontsMap, 200, "200", "Successful!", "Successfully fetch credit line fields", nil
}

// func CreateCreditLine(operation string, newCreditLine *NewCreditLine) (bool, int, string, string, string, error) {
// 	db := database.DB

// 	createQuery := `SELECT * FROM gabaykonekfunc.crdtlnecreation($1, $2)`

// 	var response map[string]any
// 	if err := db.Raw(
// 		createQuery,
// 		operation,
// 		newCreditLine,
// 	).Scan(&response).Error; err != nil {
// 		fmt.Println(err)
// 		return false, 500, "500", status.RetCode500, "An error occured while creating credit line.", err
// 	}

// 	isSuccess := sharedfunctions.GetBoolFromMap(response, "issuccess")
// 	retcode := sharedfunctions.GetStringFromMap(response, "retcode")
// 	retcodeInt := sharedfunctions.GetIntFromMap(response, "retcodeint")
// 	status := sharedfunctions.GetStringFromMap(response, "status")
// 	message := sharedfunctions.GetStringFromMap(response, "message")

// 	if !isSuccess {
// 		return isSuccess, retcodeInt, retcode, status, message, nil
// 	}

// 	return isSuccess, retcodeInt, retcode, status, message, nil
// }

func CreateCreditLine(newCreditLine map[string]any) (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw(`SELECT * FROM gabaykonekfunc.crdtlne_creation($1)`, newCreditLine).Scan(&response).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Convert JSON string fields to proper JSON if necessary
	sharedfunctions.ConvertStringToJSONMap(response)

	// Extract the JSONB field returned by the function
	response = sharedfunctions.GetMap(response, "crdtlne_creation")

	return response, nil
}

func GetCreditLine(staffID string) (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.crdtlneretrieval($1)", staffID).Scan(&response).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	creditLineList := sharedfunctions.GetMap(response, "crdtlneretrieval")

	return creditLineList, nil
}

func GetCreditLineApproved(params map[string]any) (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getapprovedcrdtlne($1)", params).Scan(&response).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	result := sharedfunctions.GetMap(response, "getapprovedcrdtlne")
	resultData := sharedfunctions.GetMap(result, "data")

	approvedList := sharedfunctions.GetListAny(resultData, "approvedCreditLine")
	for _, item := range approvedList {
		clItem := sharedfunctions.GetMap(item, "creditLineDetails")
		fields := sharedfunctions.GetMap(clItem, "creditLineFields")
		header := sharedfunctions.GetMap(fields, "headerFields")

		header["amountinwords"] = sharedfunctions.ConvertToWords(0, sharedfunctions.GetFloatFromMap(header, "amount"))

		// fmt.Println("Amount in words: ", header["amountinwords"])
	}

	// fmt.Println("Data: ", resultData)
	// fmt.Println("Approved Credit Lines: ", approvedList)

	return result, nil
}
