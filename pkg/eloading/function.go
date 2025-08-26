package eloading

import (
	"bytes"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func fetchPromoDetails(serviceReqBody map[string]any, cagabayReqBody map[string]any) (map[string]any, error) {
	url, err := sharedfunctions.GetBaseUrl(10)
	if err != nil {
		return nil, err
	}

	method := "POST"

	pHeader := sharedfunctions.GetMap(serviceReqBody, "header")
	pReqBody := sharedfunctions.GetMap(serviceReqBody, "requestBody")

	jsonBody, err := json.Marshal(pReqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling requestBody: %v", err)
	}

	payload := bytes.NewReader(jsonBody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	for key := range pHeader {
		req.Header.Add(key, sharedfunctions.GetStringFromMap(pHeader, key))
	}
	req.Header.Add("Content-Type", "application/json")

	// Send request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Debug output (optional)
	// fmt.Println(string(body))

	// Parse response JSON
	var serviceRespBody map[string]any
	if err := json.Unmarshal(body, &serviceRespBody); err != nil {
		return nil, err
	}

	result, err := serviceResponse(serviceReqBody, cagabayReqBody, serviceRespBody)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func serviceRequest(cagabayReqBody map[string]any) (map[string]any, error) {
	db := database.DB

	var serviceReqBody map[string]any
	if err := db.Raw("SELECT * FROM eloadingfunc.getrequest($1)", cagabayReqBody).Scan(&serviceReqBody).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(serviceReqBody)
	serviceReqBody = sharedfunctions.GetMap(serviceReqBody, "getrequest")

	apiResponse, err := fetchPromoDetails(serviceReqBody, cagabayReqBody)
	if err != nil {
		return nil, err
	}

	return apiResponse, nil
}

func serviceResponse(serviceReqBody map[string]any, cagabayReqBody map[string]any, serviceRespBody map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM eloadingfunc.getresponse($1, $2, $3)", serviceReqBody, cagabayReqBody, serviceRespBody).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getresponse")

	return result, nil
}
