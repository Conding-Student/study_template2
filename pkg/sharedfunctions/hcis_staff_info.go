package sharedfunctions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Hcis(staffId string) (map[string]any, error) {
	hcisBaseURL, err := GetBaseUrl(3)
	if err != nil {
		return nil, err
	}

	basicAuth, err := GetBasicAuth(1)
	if err != nil {
		return nil, err
	}

	url := hcisBaseURL + "/HCISLink/WEBAPI/ExternalService/ViewStaffInfo"

	// JSON request payload
	requestPayload := map[string]string{
		"StaffID": staffId,
	}

	// Convert the request payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("Error marshalling JSON request:", err)
		return nil, err
	}

	// Create a POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set headers for the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", basicAuth)
	// req.Header.Set("Cookie", "cookiesession1=678B29E7CDEFGHIJKLMNOPQRSTUVB249")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var hcisResponse map[string]any
	if err := json.Unmarshal(responseBody, &hcisResponse); err != nil {
		return nil, err
	}

	return hcisResponse, nil
}
