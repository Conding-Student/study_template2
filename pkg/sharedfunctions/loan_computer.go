package sharedfunctions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func LoanAmort(principal, flatRate, proRateonHalf, n, frequency int, dateReleased string, meetingDay, dueDateType, withDST, isLumpSum, intComp, gracePeriod int) (map[string]any, error) {

	loanCalculatorBaseURL, err := GetBaseUrl(7)
	if err != nil {
		log.Println(err, loanCalculatorBaseURL)
		return nil, err
	}

	url := loanCalculatorBaseURL

	payload := map[string]any{
		"principal":     principal,
		"flatRate":      flatRate,
		"proRateonHalf": proRateonHalf,
		"n":             n,
		"frequency":     frequency,
		"dateReleased":  dateReleased,
		"meetingDay":    meetingDay,
		"dueDateType":   dueDateType,
		"withDST":       withDST,
		"isLumpSum":     isLumpSum,
		"intComp":       intComp,
		"gracePeriod":   gracePeriod,
	}

	// fmt.Println(payload)

	// Encode the request body as JSON
	requestBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding request body in loan amort API:", err)
		return nil, err
	}

	// Create a POST request with the encoded request body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request in loan amort API:", err)
		return nil, err
	}

	// Set headers for the request
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request in loan amort API:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body in loan amort API:", err)
		return nil, err
	}

	// Print the response body as a string
	// fmt.Println("Response Body:", string(responseBody))
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the return code is not 200 in loan amort API")
	}

	// Check if response body is empty
	if len(responseBody) == 0 {
		return nil, fmt.Errorf("empty response body of loan amort API")
	}

	var amortResponse map[string]any
	if err := json.Unmarshal(responseBody, &amortResponse); err != nil {
		return nil, err
	}

	return amortResponse, nil
}
