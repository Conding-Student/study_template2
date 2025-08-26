package sharedfunctions

import (
	"bytes"
	"chatbot/pkg/models/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func SaveToSoteriaGoMobile(staffInfo map[string]any, username, password string, instiCode int) (string, error) {

	// fmt.Println("Soteria Stage: ", staffInfo)

	soteriaGo, err := GetBaseUrl(9)
	if err != nil {
		return "", err
	}
	url := soteriaGo + "/api/auth/signup"

	userForSignUp := map[string]any{
		"role":        []string{"cagabayUser"},
		"username":    username,
		"mobile":      GetStringFromMap(staffInfo, "mobile"),
		"password":    password,
		"birthday":    GetStringFromMap(staffInfo, "birthdate"),
		"cid":         "0",
		"instiCode":   instiCode,
		"brCode":      0,
		"unitCode":    0,
		"centerCode":  0,
		"deviceModel": GetStringFromMap(staffInfo, "deviceModel"),
		"firstName":   GetStringFromMap(staffInfo, "firstName"),
		"middleName":  GetStringFromMap(staffInfo, "middleName"),
		"lastName":    GetStringFromMap(staffInfo, "lastName"),
		"deviceId":    GetStringFromMap(staffInfo, "deviceID"),
	}

	// fmt.Println(user.DeviceID)
	// fmt.Println(user.DeviceModel)
	requestBody, err := json.Marshal(userForSignUp)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse error response: %w", err)
	}

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		// User information saved to the external API successfully
		return "", nil
	} else {
		// fmt.Println(response)
		return GetStringFromMap(response, "message"), nil
	}

}

func SaveToSoteriaGo(user *model.UserRequestHeader, username, password string, instiCode int) (string, error) {
	soteriaGo, err := GetBaseUrl(9)
	if err != nil {
		return "", err
	}
	url := soteriaGo + "/api/auth/signup"

	userForSignUp := map[string]any{
		"role":        []string{"cagabayUser"},
		"username":    username,
		"mobile":      user.ContactInfo.Mobile,
		"password":    password,
		"birthday":    user.PersonalInfo.Birthdate,
		"cid":         "0",
		"instiCode":   instiCode,
		"brCode":      0,
		"unitCode":    0,
		"centerCode":  0,
		"deviceModel": user.DeviceInfo.DeviceModel,
		"firstName":   user.PersonalInfo.FirstName,
		"middleName":  user.PersonalInfo.MiddleName,
		"lastName":    user.PersonalInfo.LastName,
		"deviceId":    user.DeviceInfo.DeviceUsed,
	}

	// fmt.Println(user.DeviceID)
	// fmt.Println(user.DeviceModel)
	requestBody, err := json.Marshal(userForSignUp)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse error response: %w", err)
	}

	// Check the response status code
	switch resp.StatusCode {
	case http.StatusOK:
		// User information saved to the external API successfully
		return "", nil
	case 400:
		// fmt.Println(response)
		return GetStringFromMap(response, "message"), nil
	}
	return "", nil
}

func SoteriaLogin(username, password, deviceId, deviceModel string) (string, string, error) {
	soteriaPy, err := GetBaseUrl(9)
	if err != nil {
		return "Problem connecting to database", "", err
	}
	url := soteriaPy + "/Soteria/api/auth/signin"

	payload := map[string]string{
		"username": username,
		"password": password,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "An error occured while connecting to server", "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "An error occured while connecting to server", "", err
	}

	// Set the Content-Type header to indicate JSON payload
	req.Header.Set("Content-Type", "application/json")
	// Set the deviceId header
	req.Header.Set("deviceId", deviceId)
	req.Header.Set("deviceModel", deviceModel)
	req.Header.Set("fcmToken", "CA-GABAY@2023")
	req.Header.Set("osVersion", "0")
	req.Header.Set("appVersion", "0")

	// Make the POST request to the authentication endpoint
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "An error occured while connecting to server", "", err
	}
	defer resp.Body.Close()

	// Print the response status code for debugging
	fmt.Println("Response Status Code:", resp.StatusCode)

	// Read the response body
	responseBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "An error occured while connecting to server", "", readErr
	}

	var response map[string]any
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "An error occured while connecting to server", "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:

		accessToken, ok := response["accessToken"].(string)
		if !ok {
			return "An error occured while connecting to server", "", errors.New("token not found in response")
		}

		return "", accessToken, nil
	case 400:

		return "Invalid credentials. Please double check your staff id/username and password", "", errors.New("error response")

	default:

		return "Unathorized request", "", errors.New("error response")

	}
}

func ChangePasswordToSoteria(username, oldpassword, newpassword string) (string, string, error) {
	soteriaPy, err := GetBaseUrl(9)
	if err != nil {
		return "Problem connecting to database", "", err
	}
	url := soteriaPy + "/Soteria/api/auth/changePassword"

	payload := map[string]string{
		"username":    username,
		"oldPassword": oldpassword,
		"password":    newpassword,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "An error occured while connecting to server", "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "An error occured while connecting to server", "", err
	}

	// Set the Content-Type header to indicate JSON payload
	req.Header.Set("Content-Type", "application/json")

	// Make the POST request to the authentication endpoint
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "An error occured while connecting to server", "", err
	}
	defer resp.Body.Close()

	// Print the response status code for debugging
	fmt.Println("Response Status Code:", resp.StatusCode)

	// Read the response body
	responseBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "An error occured while connecting to server", "", readErr
	}

	var response map[string]any
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "An error occured while connecting to server", "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:

		accessToken, ok := response["accessToken"].(string)
		if !ok {
			return "An error occured while connecting to server", "", errors.New("token not found in response")
		}

		return "", accessToken, nil
	case 400:

		return "Invalid credentials. Please double check your staff id/username and password", "", errors.New("error response")

	default:

		return "Unathorized request", "", errors.New("error response")

	}
}
