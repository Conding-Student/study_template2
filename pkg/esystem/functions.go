package esystem

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"log"
)

type AllCurrentLoan struct {
	TotalPrincipalAmount float64 `gorm:"column:total_loan_principal_amount"`
	TotalDue             float64 `gorm:"column:total_weekly_due"`
}

type GetClientAddress struct {
	Sitio    string `gorm:"column:addressdetails"`
	Barangay string `gorm:"column:barangay"`
	City     string `gorm:"column:city"`
	Province string `gorm:"column:province"`
}

type GetSavingsRequiredResponse struct {
	// BrCode          string  `gorm:"column:brcode"`
	// Cid             int     `gorm:"column:cid"`
	Savings     float64 `gorm:"column:savings"`
	SavingsType int     `gorm:"column:savingstype"`
	UpdateAsOf  string  `gorm:"column:dolasttrans"`
}

var LOSfeature = "eSystem"

func GetClientPreviousLoan(reqDetails map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getclientprevln($1)", reqDetails).Scan(&result).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getclientprevln")
	return result, nil
}

func GetClientInformationV2(clientCreds map[string]any) (map[string]any, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getclientinfo($1)", clientCreds).Scan(&response).Error; err != nil {
		log.Println(err)
		return nil, "500", err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	result := sharedfunctions.GetMap(response, "getclientinfo")
	retCode := sharedfunctions.GetStringFromMap(result, "retCode")

	return result, retCode, nil
}

func GetCoBorrowerInformation(staffid string) ([]map[string]any, error) {
	db := database.DB

	var response []map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.coborrowersdata($1)", staffid).Scan(&response).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONList(response)

	if len(response) == 0 {
		response = []map[string]any{}
	}

	return response, nil
}

func GetAllCurrentLoans(cid int, brcode string, loanProductCode int) (*AllCurrentLoan, error) {
	eSystem := database.EsystemDB
	var allCurrentLoans AllCurrentLoan

	if err := eSystem.Raw("SELECT * FROM cagabay_get_total_current_loan_details($1, $2, $3)", cid, brcode, loanProductCode).Scan(&allCurrentLoans).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	return &allCurrentLoans, nil
}

func GetClientAddresses(cid int, brcode string) (*GetClientAddress, error) {
	eSystem := database.EsystemDB
	var getClientAddress GetClientAddress

	if err := eSystem.Raw("SELECT * FROM cagabay_getclientaddress($1, $2)", cid, brcode).Scan(&getClientAddress).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	return &getClientAddress, nil
}

func GetClientSavings(cid int, brcode string) (*GetSavingsRequiredResponse, error) {
	eSystem := database.EsystemDB
	var getSavingsRequiredResponse GetSavingsRequiredResponse

	if err := eSystem.Raw("SELECT * FROM cagabay_get_client_savings_balance($1, $2)", cid, brcode).Scan(&getSavingsRequiredResponse).Error; err != nil {
		return nil, err
	}

	return &getSavingsRequiredResponse, nil
}
