package empc

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/sharedfunctions"
	"strings"

	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type EmpcLoanBalRequest struct {
	Cid         string `json:"empcCid"`
	LoanAccount string `json:"loanAccount"`
}

type EmpcLoanBalResponse struct {
	LoanProduct            string `gorm:"not null"`
	Message                string `gorm:"not null"`
	Balance                string `gorm:"not null"`
	PaidOutstandingBalance string `gorm:"not null"`
}

var EMPCFeature = "EMPC"

func CheckLoanBalance(c *fiber.Ctx) error {
	fmt.Println(divider)

	staffID := c.Params("id")
	empcLoanBalanceRequest := new(EmpcLoanBalRequest)
	if err := c.BodyParser(empcLoanBalanceRequest); err != nil {
		fmt.Println("retCode 401")
		fmt.Println("Invalid Request")
		fmt.Println("Failed to parse request", err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// var loanProducts string
	// err := db.Raw("SELECT loan_products FROM empc.empc_loans WHERE loan_id = ?", empcLoanBalanceRequest.LoanAccount).Row().Scan(&loanProducts)
	// if err != nil {
	// 	return c.Status(500).JSON(response.ResponseModel{
	// 		RetCode: "500",
	// 		Message: "Internal Server Error",
	// 		Data: errors.ErrorModel{
	// 			Message:   "Failed to retrieve loan product information",
	// 			IsSuccess: false,
	// 			Error:     err,
	// 		},
	// 	})
	// }

	loanProducts, err := getLoanProducts(empcLoanBalanceRequest.LoanAccount)
	if err != nil {
		logs.LOSLogs(c, EMPCFeature, staffID, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Problem connecting to server",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	loanResponseData, err := EmpcCheckLoanBal(empcLoanBalanceRequest.Cid, empcLoanBalanceRequest.LoanAccount)
	if err != nil {
		logs.LOSLogs(c, EMPCFeature, staffID, "404", err.Error())
		return c.Status(404).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch Empc Loan Balance",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	var empcLoanResponse map[string]any
	if err := json.Unmarshal(loanResponseData, &empcLoanResponse); err != nil {
		logs.LOSLogs(c, EMPCFeature, staffID, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to parse empc response",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	var paidOutstandingBal string
	var balance string
	//var note string
	// messageInAList := empcLoanResponse["content"].(map[string]any)["messages"].([]any)[0].(map[string]any)
	empcResp := sharedfunctions.GetMap(empcLoanResponse, "content")
	empcMess := sharedfunctions.GetListAny(empcResp, "messages")
	messageInAList := sharedfunctions.GetMapAtListIndex(empcMess, 0)

	if len(messageInAList) == 0 {
		messageInAList = fiber.Map{
			"Paid Outstanding Bal.": "0",
			"text":                  "You don't have an Existing Loan or Pending Loan Application.",
			"Balance":               "0",
		}
	}

	// messageText := messageInAList["text"].(string)
	// balanceText := messageInAList["Balance"].(string)
	// paidOutstandingBal = messageInAList["Paid Outstanding Bal."].(string)

	messageText := sharedfunctions.GetStringFromMap(messageInAList, "text")
	balanceText := sharedfunctions.GetStringFromMap(messageInAList, "Balance")
	paidOutstandingBal = sharedfunctions.GetStringFromMap(messageInAList, "Paid Outstanding Bal.")

	balanceText = strings.ReplaceAll(balanceText, "$", "₱")

	balance = balanceText

	if balance == "" || paidOutstandingBal == "" {
		paidOutstandingBal = "100%"
		balance = "₱0.00"
		if strings.Contains(messageText, "Existing Loan or Pending Loan Application.") {
			messageText = strings.Replace(messageText, "Existing Loan or Pending Loan", "Existing or Pending "+loanProducts, 1)
		}
		//note = "You can apply for a loan in EMPC."
	}

	// percent, err := strconv.Atoi(paidOutstandingBal)
	// if err != nil {
	// 	log.Println("error parsing balance.", err.Error())
	// }

	// if percent > 30 {
	// 	note = "Your loan is eligible for renewal.."
	// }

	if strings.Contains(messageText, "Existing Loan") {
		messageText = strings.Replace(messageText, "Loan", loanProducts, 1)
	}

	//fmt.Println(note)
	loanData := EmpcLoanBalResponse{
		LoanProduct:            loanProducts,
		Message:                messageText,
		Balance:                balance,
		PaidOutstandingBalance: paidOutstandingBal,
	}

	logs.LOSLogs(c, EMPCFeature, staffID, "200", ("Successfully fetch empc loan balance. " + loanProducts))
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    loanData,
	})
}
