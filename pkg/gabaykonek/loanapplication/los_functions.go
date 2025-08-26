package loans

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func AutoCancellationOfPendingApprovedLoans(c *fiber.Ctx) error {
	fmt.Println("Auto cancellation begins")
	db := database.DB

	if err := db.Exec("SELECT cron.autocancellation()").Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Auto cancellation failed!",
				IsSuccess: false,
				Error:     err,
			},
		})
	}
	fmt.Println("Auto cancellation finished")

	return c.Next()
}

func LoanCreation(loanCreation map[string]any) (map[string]any, bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.loancreationv2($1)", loanCreation).Scan(&response).Error; err != nil {
		return nil, false, 500, "500", status.RetCode500, "An error occured while creating loan application", err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	result := sharedfunctions.GetMap(response, "response")
	isSuccess := sharedfunctions.GetBoolFromMap(result, "issuccess")
	retCode := sharedfunctions.GetStringFromMap(result, "retcode")
	retCodeInt := sharedfunctions.GetIntFromMap(result, "retcode")
	status := sharedfunctions.GetStringFromMap(result, "status")
	message := sharedfunctions.GetStringFromMap(result, "message")

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nLoan Creation Successful: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return nil, isSuccess, retCodeInt, retCode, status, message, fmt.Errorf(message)
	}

	qrDetails := sharedfunctions.GetMap(result, "qrcode")

	return qrDetails, isSuccess, retCodeInt, retCode, status, message, nil
}

func LoanUpdating(loanApplication map[string]any) (bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.loanupdates($1)", loanApplication).Scan(&response).Error; err != nil {
		fmt.Println(err)
		return false, 500, "500", status.RetCode500, "An error occured while updating loan application", err
	}

	sharedfunctions.ConvertStringToJSONMap(response)

	result := sharedfunctions.GetMap(response, "response")
	retCodeInt := sharedfunctions.GetIntFromMap(result, "retcode")
	retCode := sharedfunctions.GetStringFromMap(result, "retcode")
	status := sharedfunctions.GetStringFromMap(result, "status")
	isSuccess := sharedfunctions.GetBoolFromMap(result, "issuccess")
	message := sharedfunctions.GetStringFromMap(result, "message")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nLoan update successful: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return isSuccess, retCodeInt, retCode, status, message, fmt.Errorf(message)
	}

	return isSuccess, retCodeInt, retCode, status, message, nil
}

func GetLoanLrf(amount float64, frequency, term int) (float64, error) {
	db := database.DB

	query := `SELECT * FROM gabaykonekfunc.getloanlrf($1, $2, $3)`

	var lrf float64
	if err := db.Raw(query, amount, frequency, term).Scan(&lrf).Error; err != nil {
		return 0.0, err
	}

	return lrf, nil
}

func GetLoanLoanProducDetails(lprcode int) (map[string]any, error) {
	db := database.DB
	var loanProductDetails map[string]any

	if err := db.Raw("SELECT * FROM loan_application.getloanproductspecificdetails($1)", lprcode).Scan(&loanProductDetails).Error; err != nil {
		return nil, err
	}

	return loanProductDetails, nil
}

func LoanCalculator(computeLoan map[string]any) (*LoanCalculatorPlusResponse, bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.computeloan($1)", computeLoan).Scan(&response).Error; err != nil {
		fmt.Println(err)
		return nil, false, 500, "500", status.RetCode500, "An error occured while connecting to database", err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	resultData := sharedfunctions.GetMap(response, "response")

	isSuccess := sharedfunctions.GetBoolFromMap(resultData, "issuccess")
	cStatus := sharedfunctions.GetStringFromMap(resultData, "status")
	retCode := sharedfunctions.GetStringFromMap(resultData, "retcode")
	retCodeInt := sharedfunctions.GetIntFromMap(resultData, "retcode")
	message := sharedfunctions.GetStringFromMap(resultData, "message")

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nLoan Computation Successful: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return nil, isSuccess, retCodeInt, retCode, cStatus, message, fmt.Errorf(message)
	}

	loanDetails := sharedfunctions.GetMap(resultData, "loandetails")
	computationDetails := sharedfunctions.GetMap(resultData, "computationdetails")

	loanAmortization, err := sharedfunctions.LoanAmort(
		sharedfunctions.GetIntFromMap(computationDetails, "principal"),
		sharedfunctions.GetIntFromMap(loanDetails, "flatrate"),
		sharedfunctions.GetIntFromMap(loanDetails, "prorate"),
		sharedfunctions.GetIntFromMap(computationDetails, "loanterm"),
		sharedfunctions.GetIntFromMap(computationDetails, "frequency"),
		sharedfunctions.GetStringFromMap(computationDetails, "datereleased"),
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

	startDate, err := sharedfunctions.FormatLoanComputerDate(sharedfunctions.GetStringFromMap(computationDetails, "firstpayment"))
	if err != nil {
		log.Println(err)
	}
	maturityDate, err := sharedfunctions.FormatLoanComputerDate(sharedfunctions.GetStringFromMap(computationDetails, "lastpayment"))
	if err != nil {
		log.Println(err)
	}

	amortizationList := sharedfunctions.GetListAny(loanAmortization, "Amortization")

	responses := LoanCalculatorPlusResponse{
		LoanAmount:      fmt.Sprintf("%.2f", sharedfunctions.GetFloatFromMap(computationDetails, "principal")),
		EIR:             fmt.Sprintf("%.2f %%", sharedfunctions.GetFloatFromMap(computationDetails, "eir")),
		ContractualRate: sharedfunctions.GetStringFromMap(loanAmortization, "Rate"),
		Interest:        fmt.Sprintf("%.2f", sharedfunctions.GetFloatFromMap(computationDetails, "interest")),
		InterestRate:    fmt.Sprintf("%.2f %%", sharedfunctions.GetFloatFromMap(computationDetails, "intrate")),
		LoanOutstanding: fmt.Sprintf("%.2f", sharedfunctions.GetFloatFromMap(computationDetails, "loanouts")),
		NumberOfMonths:  fmt.Sprintf("%.d", sharedfunctions.GetIntFromMap(computationDetails, "months")),
		NumberOfWeeks:   sharedfunctions.GetStringFromMap(computationDetails, "term"),
		LRF:             fmt.Sprintf("%.2f", sharedfunctions.GetFloatFromMap(computationDetails, "lrf")),
		LoanProceeds:    fmt.Sprintf("%.2f", sharedfunctions.GetFloatFromMap(computationDetails, "proceeds")),
		AmountInWords:   sharedfunctions.ConvertToWords(0, sharedfunctions.GetFloatFromMap(computationDetails, "proceeds")),
		WeeklyDue:       fmt.Sprintf("%.f", sharedfunctions.GetFloatFromMap(computationDetails, "due")),
		ReleaseDate:     loanAmortization["DateReleased2"].(string),
		StartDate:       startDate,
		MaturityDate:    maturityDate,
		Amortization:    amortizationList,
	}

	return &responses, true, 200, "200", "Successful!", "Loan succesfully computed", nil
}

func GetAllLoans(staffID string) (map[string]any, bool, int, string, string, string, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getloanapplications($1)", staffID).Scan(&response).Error; err != nil {
		fmt.Println(err)
		return nil, false, 500, "500", status.RetCode500, "An error occured while conecting to database.", err
	}

	sharedfunctions.ConvertStringToJSONMap(response)
	result := sharedfunctions.GetMap(response, "response")

	isSuccess := sharedfunctions.GetBoolFromMap(result, "issuccess")
	status := sharedfunctions.GetStringFromMap(result, "status")
	retCode := sharedfunctions.GetStringFromMap(result, "retcode")
	retCodeInt := sharedfunctions.GetIntFromMap(result, "retcode")
	message := sharedfunctions.GetStringFromMap(result, "message")

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nFetch Loan Successful: ", isSuccess)
	fmt.Println("Message: ", message)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	if !isSuccess {
		return nil, isSuccess, retCodeInt, retCode, status, message, fmt.Errorf(message)
	}

	return result, isSuccess, retCodeInt, retCode, status, message, nil
}

func GetPpiQuestionaire() (map[string]any, error) {
	db := database.DB

	var response map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getppiquestionaire()").Scan(&response).Error; err != nil {
		return nil, err
	}

	return response, nil
}
