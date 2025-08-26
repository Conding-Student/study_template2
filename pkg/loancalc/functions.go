package loancalc

import (
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
	"log"
)

func LoanCalcForBanks(computeLoan map[string]any) (*LoanResponse, bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM public.computeloan($1)", computeLoan).Scan(&response).Error; err != nil {
		fmt.Println(err)
		return nil, false, 500, "500", status.RetCode500, "An error occured while connecting to database", err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	resulData := sharedfunctions.GetMap(response, "response")

	isSuccess := sharedfunctions.GetBoolFromMap(resulData, "issuccess")
	cStatus := sharedfunctions.GetStringFromMap(resulData, "status")
	retCode := sharedfunctions.GetStringFromMap(resulData, "retcode")
	retCodeInt := sharedfunctions.GetIntFromMap(resulData, "retcode")
	message := sharedfunctions.GetStringFromMap(resulData, "message")

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nLoan Computation Successful: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return nil, isSuccess, retCodeInt, retCode, cStatus, message, fmt.Errorf("an error occured while connecting to database")
	}

	loanDetails := sharedfunctions.GetMap(resulData, "loandetails")
	computationDetails := sharedfunctions.GetMap(resulData, "computationdetails")

	loanAmortization, err := sharedfunctions.LoanAmort(
		sharedfunctions.GetIntFromMap(computationDetails, "principal"),
		sharedfunctions.GetIntFromMap(loanDetails, "flatrate"),
		sharedfunctions.GetIntFromMap(loanDetails, "prorate"),
		sharedfunctions.GetIntFromMap(computationDetails, "weeks"),
		sharedfunctions.GetIntFromMap(computationDetails, "frequency"),
		sharedfunctions.GetStringFromMap(computationDetails, "daterelease"),
		sharedfunctions.GetIntFromMap(computationDetails, "meetingday"),
		sharedfunctions.GetIntFromMap(loanDetails, "duetype"),
		sharedfunctions.GetIntFromMap(loanDetails, "withdst"),
		sharedfunctions.GetIntFromMap(computationDetails, "islumpsum"),
		sharedfunctions.GetIntFromMap(loanDetails, "intcomp"),
		sharedfunctions.GetIntFromMap(loanDetails, "graceperiod"),
	)
	if err != nil {
		return nil, false, 500, "500", status.RetCode500, "An error occured while connecting to loan amortization server", err
	}

	startDate, err := sharedfunctions.FormatLoanComputerDate(loanAmortization["DateStart"].(string))
	if err != nil {
		log.Println(err)
	}
	maturityDate, err := sharedfunctions.FormatLoanComputerDate(loanAmortization["DateMatured"].(string))
	if err != nil {
		log.Println(err)
	}

	amortizationList := sharedfunctions.GetListAny(loanAmortization, "Amortization")

	responses := LoanResponse{
		LoanAmount:       sharedfunctions.GetStringFromMap(loanAmortization, "Principal"),
		Contractual:      sharedfunctions.GetStringFromMap(loanAmortization, "Rate"),
		EIR:              sharedfunctions.GetStringFromMap(loanAmortization, "BspEffective"),
		Interest:         sharedfunctions.GetStringFromMap(loanAmortization, "Interest"),
		NumberOfMonths:   sharedfunctions.GetStringFromMap(computationDetails, "months"),
		NumberOfWeeks:    sharedfunctions.GetStringFromMap(computationDetails, "weeks"),
		LRF:              sharedfunctions.GetStringFromMap(computationDetails, "lrf"),
		DocumentaryStamp: sharedfunctions.GetStringFromMap(computationDetails, "dst"),
		LoanBalance:      sharedfunctions.GetStringFromMap(computationDetails, "loanbalance"),
		TotalDeduction:   sharedfunctions.GetStringFromMap(computationDetails, "totaldeduction"),
		LoanProceeds:     sharedfunctions.GetStringFromMap(computationDetails, "loanproceeds"),
		WeeklyDue:        sharedfunctions.GetStringFromMap(loanAmortization, "DueAmt"),
		ReleaseDate:      sharedfunctions.GetStringFromMap(loanAmortization, "DateReleased2"),
		StartDate:        startDate,
		MaturityDate:     maturityDate,
		AmountInWords:    sharedfunctions.ConvertToWords(1, sharedfunctions.GetFloatFromMap(computationDetails, "loanproceeds")),
		Amortization:     amortizationList,
	}

	return &responses, true, 200, "200", "Successful!", "Loan succesfully computed", nil
}

func RBICalcForBanks(computeLoan map[string]any) (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM rbifunc.computeloan($1)", computeLoan).Scan(&response).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	resulData := sharedfunctions.GetMap(response, "computeloan")
	data := sharedfunctions.GetMap(resulData, "data")
	data["amountInWords"] = sharedfunctions.ConvertToWords(1, sharedfunctions.GetFloatFromMap(data, "loanProceeds"))

	return resulData, nil
}
