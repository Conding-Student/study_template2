package empc

import (
	"bytes"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getLoanProducts(loanProductCode string) (string, error) {
	db := database.DB
	var loanProducts string
	err := db.Raw("SELECT loan_products FROM empc.empc_loans WHERE loan_id = ?", loanProductCode).Row().Scan(&loanProducts)
	if err != nil {
		return "", err
	}
	return loanProducts, nil
}

func GetSOA(params map[string]any) (map[string]any, error) {
	db := database.DB

	empcBaseUrl, err := sharedfunctions.GetBaseUrl(5)
	if err != nil {
		return nil, err
	}
	// url := "https://test-empc-uat.cardmri.com:9446/EMPCChatbot/webapi/EMPCService/getSOA"
	url := empcBaseUrl + "/EMPCChatbot/webapi/EMPCService/getSOA"

	cid := sharedfunctions.GetStringFromMap(params, "cid")
	staffid := sharedfunctions.GetStringFromMap(params, "staffid")

	requestPayload := map[string]string{
		"cid": cid,
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("Error marshalling JSON request:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "cookiesession1=678B28F7EGHIJKLMNOPQRSTUVWXYDCF8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var statementOfAccount map[string]any
	if err := json.Unmarshal(responseBody, &statementOfAccount); err != nil || len(statementOfAccount) == 0 {
		// fallback to stored SOA
		fallbackPayload := map[string]any{
			"staffid":   staffid,
			"querytype": 1,
		}
		var storedSOA map[string]any
		if err := db.Raw("SELECT * FROM empcfunc.soa($1)", fallbackPayload).Scan(&storedSOA).Error; err != nil {
			return nil, err
		}
		sharedfunctions.ConvertStringToJSONMap(storedSOA)
		return sharedfunctions.GetMap(storedSOA, "soa"), nil
	}

	// Primary source is successful
	statementOfAccount["querytype"] = 0
	statementOfAccount["staffid"] = staffid

	var storedSOA map[string]any
	if err := db.Raw("SELECT * FROM empcfunc.soa($1)", statementOfAccount).Scan(&storedSOA).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(storedSOA)
	return sharedfunctions.GetMap(storedSOA, "soa"), nil
}

func GetEMPCLoan(params map[string]any) (map[string]any, error) {
	db := database.DB

	empcBaseUrl, err := sharedfunctions.GetBaseUrl(5)
	if err != nil {
		return nil, err
	}
	url := empcBaseUrl + "/EMPCChatbot/webapi/EMPCService/getAmortization"

	accType := sharedfunctions.GetIntFromMap(params, "acctType")
	staffid := sharedfunctions.GetStringFromMap(params, "staffid")

	requestPayload := map[string]any{
		"hcisId":   staffid,
		"acctType": accType,
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("Error marshalling JSON request:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "cookiesession1=678B28F7EGHIJKLMNOPQRSTUVWXYDCF8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var amortization []map[string]any
	if err := json.Unmarshal(responseBody, &amortization); err != nil || len(amortization) == 0 {
		fmt.Println("Error response from external API: ", err)
		// fallback to stored SOA
		fallbackPayload := map[string]any{
			"staffid":   staffid,
			"querytype": 1,
			"acctType":  accType,
		}
		var storedSOA map[string]any
		if err := db.Raw("SELECT * FROM empcfunc.amortization($1)", fallbackPayload).Scan(&storedSOA).Error; err != nil {
			return nil, err
		}
		sharedfunctions.ConvertStringToJSONMap(storedSOA)
		return sharedfunctions.GetMap(storedSOA, "amortization"), nil
	}

	// Primary source is successful
	amortSched := map[string]any{
		"querytype":    0,
		"staffid":      staffid,
		"amortization": amortization,
		"acctType":     accType,
	}

	var storedAmort map[string]any
	if err := db.Raw("SELECT * FROM empcfunc.amortization($1)", amortSched).Scan(&storedAmort).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(storedAmort)
	return sharedfunctions.GetMap(storedAmort, "amortization"), nil
}

func EmpcempcStaffInfo(staffId string) ([]byte, error) {
	empcBaseUrl, err := sharedfunctions.GetBaseUrl(5)
	if err != nil {
		return nil, err
	}
	// url := "https://test-empc-uat.cardmri.com:9446/EMPCChatbot/webapi/EMPCService/getClientInfo"
	url := empcBaseUrl + "/EMPCChatbot/webapi/EMPCService/getClientInfo"

	// JSON request payload
	requestPayload := map[string]string{
		"staffID": staffId,
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
	req.Header.Set("Cookie", "cookiesession1=678B28F7EGHIJKLMNOPQRSTUVWXYDCF8") // Add your access token or any other required headers

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

	return responseBody, nil
}

func EmpcCheckLoanBal(cid string, loanAccount string) ([]byte, error) {
	empcBaseUrl, err := sharedfunctions.GetBaseUrl(5)
	if err != nil {
		return nil, err
	}

	// url := "https://test-empc-uat.cardmri.com:9446/EMPCChatbot/webapi/EMPCService/getMajorStatus"
	url := empcBaseUrl + "/EMPCChatbot/webapi/EMPCService/getMajorStatus"

	// JSON request payload
	requestPayload := map[string]string{
		"cid":      cid,
		"acctType": loanAccount,
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
	req.Header.Set("Cookie", "cookiesession1=678B28F7EGHIJKLMNOPQRSTUVWXYDCF8") // Add your access token or any other required headers

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

	return responseBody, nil
}
