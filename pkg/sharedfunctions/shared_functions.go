package sharedfunctions

import (
	"chatbot/pkg/utils/go-utils/database"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/divan/num2words"
)

func LocalTime() (string, error) {
	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		log.Println(err)
		return "Failed to get local time", err
	}

	manilaTime := time.Now().In(loc)
	localTime := manilaTime.Format("2006-01-02 15:04:05.999999-07:00")

	return localTime, nil
}

func LocalDateTimeWithoutSeconds() (string, error) {
	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		log.Println(err)
		return "Failed to get local time", err
	}

	manilaTime := time.Now().In(loc)
	localTime := manilaTime.Format("2006-01-02 15:04:05")

	return localTime, nil
}

func LocalDateWithoutTime() (string, error) {
	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		log.Println(err)
		return "Failed to get local time", err
	}

	manilaTime := time.Now().In(loc)
	localTime := manilaTime.Format("2006-01-02")

	return localTime, nil
}

func FormatToDateOnly(date string) (string, error) {

	formats := []string{
		time.RFC3339,  // "2006-01-02T15:04:05Z07:00"
		"2006-01-02",  // "2025-03-07"
		"02-Jan-2006", // "07-Mar-2025"
		"02/01/2006",  // "07/03/2025"
		"01/02/2006",  // "03/07/2025"
		"2006/01/02",  // "2025/03/07"
	}

	for _, layout := range formats {
		if t, err := time.Parse(layout, date); err == nil {
			return t.Format("2006-01-02"), nil
		}
	}

	return "", fmt.Errorf("unrecognized date format: %s", date)
}

func FormatLoanComputerDate(dateStr string) (string, error) {
	formats := []string{
		time.RFC3339,  // "2006-01-02T15:04:05Z07:00"
		"2006-01-02",  // "2025-03-07"
		"02-Jan-2006", // "07-Mar-2025"
		"02/01/2006",  // "07/03/2025"
		"01/02/2006",  // "03/07/2025"
		"2006/01/02",  // "2025/03/07"
	}

	for _, layout := range formats {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("02-Jan-2006"), nil
		}
	}

	return "", fmt.Errorf("unrecognized date format: %s", dateStr)
}

func GetBaseUrl(id int) (string, error) {
	db := database.DB

	var baseUrl string
	if err := db.Raw("SELECT dns FROM baseurl.url WHERE id = ?", id).Scan(&baseUrl).Error; err != nil {
		return "Problem connecting to server.", err
	}

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nExternal API from Get Base URL function: ", baseUrl)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	return baseUrl, nil
}

func GetBasicAuth(id int) (string, error) {
	db := database.DB

	var basicAuth string
	if err := db.Raw("SELECT basicauthentication FROM authentication.basicauth WHERE id = ?", id).Scan(&basicAuth).Error; err != nil {
		return "Problem connecting to server.", err
	}

	// fmt.Println("Basic Auth: ", basicAuth)
	return basicAuth, nil
}

func ConvertStringToFloatingNumber(interest string) (float64, error) {
	floatInterest, err := strconv.ParseFloat(interest, 64)
	if err == nil {
		return floatInterest, nil
	}

	return 0, fmt.Errorf("failed to convert interest string to floating numbers")
}

func ConvertToWords(convertion int, amount float64) string {
	if convertion == 0 {
		pesos := int(amount)
		return fmt.Sprintf("%s pesos only", num2words.Convert(pesos))
	} else {
		pesos := int(amount)
		cents := int((amount - float64(pesos)) * 100)
		return fmt.Sprintf("%s pesos and %s centavos only",
			num2words.Convert(pesos),
			num2words.Convert(cents),
		)
	}
}

func DesignationToInt(designation string) (int, error) {
	switch designation {
	case "Account Officer":
		return 0, nil
	case "Unit Manager":
		return 1, nil
	case "Acting Unit Manager":
		return 1, nil
	case "Area Manager":
		return 2, nil
	case "Regional Director":
		return 3, nil
	default:
		return -1, fmt.Errorf("unknown designation: %s", designation) // Handle unknown designations
	}
}

func GetLoanProductDetails(operation, lprcode int) (map[string]any, bool, string, string, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM loan_application.getloanproductdetails($1, $2)", operation, lprcode).Scan(&result).Error; err != nil {
		return nil, false, "", "", err
	}

	isSuccess, ok := result["issuccess"].(bool)
	if !ok {
		isSuccess = false
	}

	responseMessage, ok := result["message"].(string)
	if !ok {
		responseMessage = ""
	}

	errDetails, ok := result["errdetails"].(string)
	if !ok {
		errDetails = ""
	}

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nLoan product details fetch successful: ", isSuccess)
	fmt.Println("Message: ", responseMessage)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	return result, isSuccess, responseMessage, errDetails, nil
}

func AccountDetailsForSoteria(staffid string) (bool, string, string, string, string, string, string, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM public.accounlogininfo($1)", staffid).Scan(&result).Error; err != nil {
		return false, "", "", "", "", "", "Problem connecting to database", err
	}

	isSuccess, ok := result["issuccess"].(bool)
	if !ok {
		isSuccess = false
	}

	staffID, ok := result["staffid"].(string)
	if !ok {
		staffID = ""
	}

	username, ok := result["username"].(string)
	if !ok {
		username = ""
	}

	devicemodel, ok := result["devicemodel"].(string)
	if !ok {
		devicemodel = ""
	}

	deviceid, ok := result["deviceid"].(string)
	if !ok {
		deviceid = ""
	}

	rolename, ok := result["rolename"].(string)
	if !ok {
		rolename = ""
	}

	errDetails, ok := result["errdetails"].(string)
	if !ok {
		errDetails = ""
	}

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nStaffID/Username: ", staffID)
	fmt.Println("Successful: ", isSuccess)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	return isSuccess, staffID, username, devicemodel, deviceid, rolename, errDetails, nil
}

func GetStringFromMap(infoMap map[string]any, key string) string {
	if value, exists := infoMap[key]; exists {
		switch v := value.(type) {
		case string:
			return v
		case float64:
			// Remove decimals if it's a whole number
			if v == float64(int64(v)) {
				return strconv.FormatInt(int64(v), 10)
			}
			return strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			return strconv.Itoa(v)
		case int64:
			return strconv.FormatInt(v, 10)
		case fmt.Stringer:
			return v.String()
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func GetBoolFromMap(infoMap map[string]any, key string) bool {
	value, exists := infoMap[key].(bool)
	if !exists {
		return false
	}
	return value
}

func GetIntFromMap(infoMap map[string]any, key string) int {
	if value, exists := infoMap[key]; exists {
		switch v := value.(type) {
		case int:
			return v
		case int32:
			return int(v)
		case int64:
			return int(v)
		case float64:
			return int(v)
		case string:
			if intValue, err := strconv.Atoi(v); err == nil {
				return intValue
			}
		default:
			fmt.Printf("Unexpected type for key %s: %T\n", key, value)
		}
	}
	return 0
}

func GetFloatFromMap(infoMap map[string]any, key string) float64 {
	if value, exists := infoMap[key]; exists {
		switch v := value.(type) {
		case int:
			return float64(v)
		case int32:
			return float64(v)
		case int64:
			return float64(v)
		case float64:
			return v
		case string:
			if floatValue, err := strconv.ParseFloat(v, 64); err == nil {
				return floatValue
			}
		default:
			fmt.Printf("Unexpected type for key %s: %T\n", key, value)
		}
	}
	return 0
}

func GetDateFromMap(infoMap map[string]any, key string) string {
	if value, exists := infoMap[key]; exists {
		switch v := value.(type) {
		case string:
			layouts := []string{
				"2006-01-02T15:04:05Z07:00",
				"2006-01-02",
				"02-Jan-2006",
			}
			for _, layout := range layouts {
				if parsedTime, err := time.Parse(layout, v); err == nil {
					return parsedTime.Format("2006-01-02") // Return date-only format
				}
			}
			fmt.Printf("Could not parse date for key %s: %s\n", key, v)
		case time.Time:
			return v.Format("2006-01-02")
		case int64:
			return time.Unix(v, 0).Format("2006-01-02")
		case float64:
			return time.Unix(int64(v), 0).Format("2006-01-02")
		default:
			fmt.Printf("Unexpected type for key %s: %T\n", key, value)
		}
	}
	return ""
}

func RoundUpToTenth(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Ceil(val*ratio) / ratio
}

func NormalizePhoneNumber(phone string) string {

	if len(phone) == 11 && strings.HasPrefix(phone, "0") {
		phone = phone[1:]
	} else if len(phone) == 12 && strings.HasPrefix(phone, "63") {
		phone = phone[2:]
	} else if len(phone) == 13 && strings.HasPrefix(phone, "+63") {
		phone = phone[3:]
	}

	return phone
}

func GenerateSaveTempPass(operation int, staffid, email, mobile, encryptedPassword string) (bool, string, string, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM public.gentemppass($1, $2, $3, $4, $5)", operation, staffid, email, mobile, encryptedPassword).Scan(&result).Error; err != nil {
		return false, "Problem connecting to database.", "", err
	}

	isSuccess := GetBoolFromMap(result, "issuccess")
	resultMessage := GetStringFromMap(result, "message")
	tempPass := GetStringFromMap(result, "temporarypassword")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nSuccessful?:", isSuccess)
	fmt.Println("Message:", resultMessage)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	return isSuccess, resultMessage, tempPass, nil
}

func GenOTP(staffid string, otptype int) (bool, string, string, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM public.genotp($1, $2)", staffid, otptype).Scan(&result).Error; err != nil {
		return false, "Failed to generate One-Time Password.", "", err
	}

	isSuccess := GetBoolFromMap(result, "issuccess")
	resultMessage := GetStringFromMap(result, "message")
	otp := GetStringFromMap(result, "otp")

	if !isSuccess {
		return isSuccess, resultMessage, otp, fmt.Errorf("failed to generate One-Time Password")
	}

	return isSuccess, resultMessage, otp, nil
}

func GetMap(data map[string]any, key string) map[string]any {
	if val, ok := data[key].(map[string]any); ok {
		return val
	}
	return map[string]any{}
}

func GetList(data map[string]any, key string) []map[string]any {
	if val, ok := data[key]; ok {
		if list, ok := val.([]map[string]any); ok {
			return list
		}
	}
	return []map[string]any{}
}

func GetListString(data map[string]any, key string) []string {
	if val, ok := data[key]; ok {
		if interfaceList, ok := val.([]any); ok {
			strList := make([]string, 0)
			for _, v := range interfaceList {
				if str, ok := v.(string); ok {
					strList = append(strList, str)
				}
			}
			return strList
		}
	}
	return []string{}
}

func GetListAny(data map[string]any, key string) []map[string]any {
	if val, ok := data[key]; ok {
		if list, ok := val.([]any); ok {
			result := make([]map[string]any, 0, len(list))
			for _, item := range list {
				if itemMap, ok := item.(map[string]any); ok {
					result = append(result, itemMap)
				}
			}
			return result
		}
	}
	return []map[string]any{}
}

func ConvertStringToJSONMap(data map[string]any) {
	for key, value := range data {
		if str, ok := value.(string); ok {
			var jsonObject map[string]any
			if err := json.Unmarshal([]byte(str), &jsonObject); err == nil {
				data[key] = jsonObject
				continue
			}
			var jsonArray []map[string]any
			if err := json.Unmarshal([]byte(str), &jsonArray); err == nil {
				data[key] = jsonArray
			}
		}
	}
}

func ConvertStringToJSONList(data []map[string]any) {
	for i := range data {
		for key, value := range data[i] {
			if str, ok := value.(string); ok {
				var jsonValue map[string]any
				if err := json.Unmarshal([]byte(str), &jsonValue); err == nil {
					data[i][key] = jsonValue
				}
			}
			if str, ok := value.(string); ok {
				var jsonValue []map[string]any
				if err := json.Unmarshal([]byte(str), &jsonValue); err == nil {
					data[i][key] = jsonValue
				}
			}
		}
	}

}

func GetMapAtListIndex(data []map[string]any, index int) map[string]any {
	if index < 0 || index >= len(data) {
		return map[string]any{}
	}
	return data[index]
}

func GetMapAtListAny(data []any, index int) map[string]any {
	if index < 0 || index >= len(data) {
		return map[string]any{}
	}

	if result, ok := data[index].(map[string]any); ok {
		return result
	}

	return map[string]any{}
}
